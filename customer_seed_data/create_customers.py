import pandas as pd
import copy
import random

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
    age_weights = age.drop(age_drop_labels, axis=1).div(age['pop'], axis=0)

    gender_labels = ['male', 'female']
    gender_weights = age[gender_labels].div(age['pop'], axis=0)

    breakpoint()
    codes = random.choices(income.index.astype('int'), 
                           weights=code_weights.values.astype('float64'), 
                           k=N)
    income_brackets = [random.choices(income_labels, 
                            income_weights.loc[code]) for code in codes]
    ages = [random.choices() for code in codes]

    breakpoint()


    return 

if __name__ == "__main__":
    income, age, education, occupation = import_all()

    _ = create_customers(100, income, age, education, occupation)
    