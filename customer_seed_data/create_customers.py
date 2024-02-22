import copy
import random
import functools
import pandas as pd
import numpy as np
import time

from concurrent.futures import ProcessPoolExecutor
from import_seed import import_all


# TODO: There is too much parameter passthrough, maybe refactor
# the functions into a class.

def create_customers(N: int,
                     income: pd.DataFrame,
                     age: pd.DataFrame,
                     education: pd.DataFrame,
                     occupation: pd.DataFrame,
                     max_workers: int = 2) -> pd.DataFrame:
    """Create customers based on the provided income, age, education, and occupation statistics.

    N:          Number of customers to create. 
    income:     DataFrame containing income data as returned by 'import_income'.
    age:        DataFrame containing age data as returned by 'import_age'.
    education:  DataFrame containing education data as returned by 'import_education'.
    occupation: DataFrame containing occupation data as returned by 'import_occupation'.
    max_workers: Maximum amount of multiprocessing workers. Passed to ProcesspoolExecutor.

    returns:    DataFrame containing customer information.
    """

    # Protect original tables from getting mangled. 
    income = copy.copy(income)
    age = copy.copy(age)
    education = copy.copy(education)
    occupation = copy.copy(occupation)

    # Create weights for random.choices for picking correct distribution.
    # Divide amount of people in each category by the relevant total amount of people.
    code_labels = ['over_18']
    income_labels = ['low', 'middle', 'upper']
    occupation_labels = ['employed', 'unemployed', 'students', 'retired', 'other']
    age_labels = ['18_19', '20_24', '25_29', '30_34', '35_39',
                  '40_44', '45_49', '50_54', '55_59', '60_64',
                  '65_69', '70_74', '75_79', '80_84', '85_90']
    gender_labels = ['male', 'female']
    education_labels = ['grade', 'high_school', 'vocational', 
                        'lower_uni', 'higher_uni']

    code_weights = income[code_labels].div(income['over_18'].sum(), axis=0)
    income_weights = income[income_labels].div(income['over_18'], axis=0)
    age_weights = age[age_labels].div(age['pop'], axis=0)
    gender_weights = age[gender_labels].div(age['pop'], axis=0)
    occupation_weights = occupation[occupation_labels].div(occupation['pop'], axis=0)
    education_weights = education[education_labels].div(education['over_18'], axis=0)

    # First we pick area code for each customer, then pick the rest of the 
    # parameters based on the statistics for that area code.
    codes = random.choices(income.index.astype('int'), 
                           weights=code_weights.values.astype('float64'), 
                           k=N)

    # TODO: Make separate function for picking all the customer parameters.
    # This way we can better control the results (don't want 18 years 
    # old retirees), and easier to use in multiprocessing.map if we need
    # to create large amount of customers.

    executable = functools.partial(get_customer_parameters,
                                  income_labels = income_labels, 
                                  income_weights = income_weights,
                                  age_labels = age_labels, 
                                  age_weights = age_weights,
                                  gender_labels = gender_labels, 
                                  gender_weights = gender_weights,
                                  education_labels = education_labels, 
                                  education_weights = education_weights,
                                  occupation_labels = occupation_labels, 
                                  occupation_weights = occupation_weights,
                                  income = income
                                 )

    code_input = list(split(codes, max_workers))
    with ProcessPoolExecutor(max_workers=max_workers) as exec:
        results = exec.map(executable, code_input) 
    # breakpoint()

    column_labels = ['code', 'age', 'gender', 'occupation',
                     'education', 'income_bracket', 'income', 'active_hour']
    results_array = np.concatenate( [np.array(x) for x in list(results) ], axis=0)
    res = pd.DataFrame(results_array,
                       columns=column_labels)

    # res = pd.concat(results)

    # res = pd.DataFrame( {
    #     'code': pd.Series(codes, dtype='int'),
    #     'age': pd.Series(ages, dtype='int'),
    #     'gender': pd.Series(genders, dtype='str'),
    #     'occupation': pd.Series(occupations, dtype='str'),
    #     'education': pd.Series(educations, dtype='str'),
    #     'income_bracket': pd.Series(income_brackets, dtype='str'),
    #     'income': pd.Series(incomes, dtype='float'),
    #     'active_hour': pd.Series(most_active, dtype='float')
    # })

    return res

# TODO: Horrible amount of parameters for single function.
# TODO: Change to take single code and output single customer.
def get_customer_parameters(codes: list[int],
                            income_labels: list[str], income_weights: list[float],
                            age_labels: list[str], age_weights: list[float],
                            gender_labels: list[str], gender_weights: list[float],
                            education_labels: list[str], education_weights: list[float],
                            occupation_labels: list[str], occupation_weights: list[float],
                            income: pd.DataFrame
                            ) -> pd.DataFrame:
    # NOTE: random.choices always returns list even when retrieving single value.
    # Hence the [0].

    res = list()
    for code in codes:
        income_bracket = random.choices(income_labels, income_weights.loc[code])[0]
        age = get_age(random.choices(age_labels, age_weights.loc[code])[0])
        gender = random.choices(gender_labels, gender_weights.loc[code])[0]
        education = random.choices(education_labels, education_weights.loc[code])[0]
        occupation = random.choices(occupation_labels, occupation_weights.loc[code])[0]

        # Modulo forces values to range 0-24.
        most_active = np.random.normal(15.0, 6) % 24

        actual_income =  get_income(income.loc[code], 
                                    income_bracket,
                                    education,
                                    occupation)

        res.append((code, age, gender, 
                   occupation, education, 
                   income_bracket, actual_income, 
                   most_active))

    return res


def modify_education_weights():
    """Modify the weights for education classes based on person age.
    Young people should have lower chance of having higher education
    than 30+ people.
    """
    ...

def modify_occupation_weights():
    """Modify the weights for occupation classes based on person age and education.
    Young people should have higher chance of being students or unemployed, older people should
    have higher chance of being retired.
    Educated people should have higher chance of being employed, and vice versa. 
    """
    ...


def get_shopping_categories():
    """Create shopping categories for customers. Either completely random or make some crude stereotypes..."""
    ...


def get_age(age: str) -> int:
    """Get random age from age-label in the form of '20_25'."""

    begin, end = [int(val) for val in age.split(sep='_')]
    res = np.random.randint(begin, end+1)
    return res


def get_income(income: pd.DataFrame, 
               income_bracket: str,
               education: str, 
               occupation: str) -> float:
    """Modify income value depending on the education, occupation, 
    and income_bracket. 
    """
    if income_bracket == 'low':
        res = income['median'] * 0.80
    elif income_bracket == 'middle':
        res = income['median']
    elif income_bracket == 'upper':
        res = income['avg']

    # These values are completely made up.
    education_modifier = {'grade': 0.65,
                          'high_school': 0.80,
                          'vocational': 1.00,
                          'lower_uni': 1.10,
                          'higher_uni': 1.20}
    occupation_modifier = {'employed': 1.25,
                           'unemployed': 0.60,
                           'students': 0.70,
                           'retired': 0.85,
                           'other': 1.00}

    res *= education_modifier[education]
    res *= occupation_modifier[occupation]

    return res


def timer(f, *args, **kwargs): 
    start = time.monotonic()
    res = f(*args, **kwargs)
    dt = time.monotonic() - start
    print(f"{dt} s")

    return res


def split(a, n):
    """Split list to n parts. Stolen from stackoverflow."""

    k, m = np.divmod(len(a), n)
    return (a[i*k+min(i, m):(i+1)*k+min(i+1, m)] for i in range(n))


if __name__ == "__main__":
    income, age, education, occupation = import_all()

    customers = timer(create_customers, 5000, income, age, education, occupation, 
                      max_workers=4)
    customers = timer(create_customers, 10000, income, age, education, occupation, 
                      max_workers=4)
    customers = timer(create_customers, 20000, income, age, education, occupation, 
                      max_workers=4)
    customers = timer(create_customers, 40000, income, age, education, occupation, 
                      max_workers=4)
    customers = timer(create_customers, 80000, income, age, education, occupation, 
                      max_workers=4)
