import copy
import random
import pandas as pd
import numpy as np

from import_seed import import_all


def create_customers(N: int,
                     income: pd.DataFrame,
                     age: pd.DataFrame,
                     education: pd.DataFrame,
                     occupation: pd.DataFrame) -> pd.DataFrame:
    """Create customers based on the provided income, 
    age, education, and occupation statistics.
    

    N: number of customers to create. 
    income: DataFrame containing income data as returned by 'import_income'.
    age: DataFrame containing age data as returned by 'import_age'.
    education: DataFrame containing education data as returned by 'import_education'.
    occupation: DataFrame containing occupation data as returned by 'import_occupation'.

    -> DataFrame containing customer information.
    """
    # TODO: Very long function...

    # Protect original tables from getting mangled. 
    income = copy.copy(income)
    age = copy.copy(age)
    education = copy.copy(education)
    occupation = copy.copy(occupation)

    # Create weights for random.choices for picking correct distribution.
    # Divide amount of people in each category by the relevant total amount of people.
    #
    # NOTE: Depending on which is easier, we either select columns we want or
    # drop columns which we do not want.
    code_labels = ['over_18']
    code_weights = income[code_labels].div(income['over_18'].sum(), axis=0)

    income_labels = ['low', 'middle', 'upper']
    income_weights = income[income_labels].div(income['over_18'], axis=0)

    age_drop_labels = ['code', 'area', 'male', 'female', 'avg', 'pop']
    age_labels = age.drop(age_drop_labels, axis=1).columns
    age_weights = age.drop(age_drop_labels, axis=1).div(age['pop'], axis=0)

    gender_labels = ['male', 'female']
    gender_weights = age[gender_labels].div(age['pop'], axis=0)

    occupation_labels = ['employed', 'unemployed', 'students', 'retired', 'other']
    occupation_weights = occupation[occupation_labels].div(occupation['pop'], axis=0)


    codes = random.choices(income.index.astype('int'), 
                           weights=code_weights.values.astype('float64'), 
                           k=N)

    # TODO: Make separate function for picking all the customer parameters.
    # This way we can better control the results (don't want 18 years 
    # old retirees), and easier to use in multiprocessing.map if we need
    # to create large amount of customers.

    # NOTE: random.choices always returns list even when retrieving single value.
    # Hence the [0].
    income_brackets = [random.choices(income_labels, income_weights.loc[code])[0]
                       for code in codes]

    ages = [get_age(random.choices(age_labels, age_weights.loc[code])[0])
            for code in codes]

    genders = [random.choices(gender_labels, gender_weights.loc[code])[0]
               for code in codes]

    occupations = [random.choices(occupation_labels, occupation_weights.loc[code])[0]
                   for code in codes]

    res = pd.DataFrame( {
        'code': pd.Series(codes, dtype='int'),
        'age': pd.Series(ages, dtype='int'),
        'gender': pd.Series(genders, dtype='str'),
        'occupation': pd.Series(occupations, dtype='str'),
        'income_bracket': pd.Series(income_brackets, dtype='str')
    })

    return res


def get_age(age: str) -> int:
    """Get random age from age-label in the form of '20_25'."""

    begin, end = [int(val) for val in age.split(sep='_')]
    res = np.random.randint(begin, end+1)
    return res

if __name__ == "__main__":
    income, age, education, occupation = import_all()

    customers = create_customers(100, income, age, education, occupation)
    