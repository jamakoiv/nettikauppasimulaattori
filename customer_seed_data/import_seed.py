import pandas as pd
from pathlib import Path

from typing import Tuple

def import_income(path: str | Path) -> pd.DataFrame:
    """Import income-stats as pandas DataFrame.

    Column-labels:
        'code': Area code
        'area': Area name
        'over_18':   N Persons over 18
        'avg':  Average income €
        'median':  Median income €
        'low':  N in lowest income bracket
        'middle':  N in middle income bracket
        'upper':  N in upper income bracket
        'cumulative':  Cumulative purchasing power (?)
    """

    names = ['code', 'area', 'over_18', 'avg', 'median', 
             'low', 'middle', 'upper', 'cumulative']
    dtypes = {'code': int,      
              'area': str,      
              'over_18': 'Int64',    
              'avg': 'Int64',   
              'median': 'Int64',   
              'low': 'Int64',   
              'middle': 'Int64',   
              'upper': 'Int64',   
              'cumulative': 'Int64'}   

    na_values = ['...', ' ...']

    df = pd.read_csv(path, 
                     header=None, 
                     comment='#',
                     sep=';',
                     names=names,
                     dtype=dtypes,
                     na_values=na_values )

    return df

def import_age(path: str | Path) -> pd.DataFrame:
    """Import age-statistics as pandas DataFrame.
    
    Columns-labels:
        'code':     Area code
        'area':     Area name
        'pop':      Total population
        'male':     Male population
        'female':   Female population
        'avg':      Average age
        '18_19':    Pop. between ages 18-19
        ....
        '80_84':    Pop. between ages 80-84
        '85':       Pop. above age 85
    """

    # NOTE: Label for over 85 years old is written as '85_90' 
    # to force common pattern for all names.
    names = ['code', 'area', 'pop', 'male', 'female', 'avg', '18_19', '20_24', 
             '25_29', '30_34', '35_39', '40_44', '45_49', '50_54', '55_59',
             '60_64', '65_69', '70_74', '75_79', '80_84', '85_90']
    dtypes = {'code': int,      
              'area': str,      
              'pop': 'Int64',    
              'male': 'Int64',   
              'female': 'Int64',   
              'avg': 'Int64',   
              '18_19': 'Int64',   
              '20_24': 'Int64',   
              '25_29': 'Int64',   
              '30_34': 'Int64',   
              '35_39': 'Int64',   
              '40_44': 'Int64',   
              '45_49': 'Int64',   
              '50_54': 'Int64',   
              '55_59': 'Int64',   
              '60_64': 'Int64',   
              '65_69': 'Int64',   
              '70_74': 'Int64',   
              '75_79': 'Int64',   
              '80_84': 'Int64',   
              '85_90': 'Int64'}

    na_values = ['...', ' ...']

    df = pd.read_csv(path, 
                     header=None, 
                     comment='#',
                     sep=';',
                     names=names,
                     dtype=dtypes,
                     na_values=na_values )
    return df


def import_occupation(path: str | Path) -> pd.DataFrame:
    """Import occupation-statistics as pandas DataFrame.
    
    Columns-labels:
        code:     Area code
        area:     Area name
        pop:      Total population
        employed: Working population
        unemployed: Unemployed population
        kids:     Population under age of 18
        students: Student population
        retired:  Retirees
        other:    Probably lizard-people (idk.)
    """

    names = ['code', 'area', 'pop', 
             'employed', 'unemployed',
             'kids', 'students', 'retired', 'other']
    dtypes = {
            'code':     int, 
            'area':     str, 
            'pop':      'Int64',
            'employed': 'Int64', 
            'unemployed':'Int64',
            'kids':     'Int64',
            'students': 'Int64',
            'retired':  'Int64',
            'other':    'Int64' }

    na_values = ['...', ' ...']

    df = pd.read_csv(path, 
                     header=None, 
                     comment='#',
                     sep=';',
                     names=names,
                     dtype=dtypes,
                     na_values=na_values )
    return df


def import_education(path: str | Path) -> pd.DataFrame:
    """Import education-statistics as pandas DataFrame.
    
    Columns-labels:
        code:   Area code.
        area:   Area name.
        over_18:     Population above age of 18.     
        grade:  Finished grade school.
        educated: Finished education higher than grade school.
        high_school: Finished high school.
        vocational: Finished vocational education.
        lower_uni: Finished lower university class education.
        higher_uni: Finished higher university class education.
    """

    names =['code', 'area', 'over_18', 'grade', 'educated', 
            'high_school', 'vocational',
            'lower_uni', 'higher_uni']
    dtypes = {
            'code': int,
            'area': str, 
            'over_18': 'Int64', 
            'grade': 'Int64',
            'educated': 'Int64',
            'high_school': 'Int64',
            'vocational': 'Int64',
            'lower_uni': 'Int64',
            'higher_uni': 'Int64' }

    na_values = ['...', ' ...']

    df = pd.read_csv(path, 
                     header=None,
                     comment='#',
                     sep=';',
                     names=names,
                     dtype=dtypes,
                     na_values=na_values )
    return df


def import_all(dropna=True, index_code=True) -> Tuple[pd.DataFrame,
                                                      pd.DataFrame,
                                                      pd.DataFrame,
                                                      pd.DataFrame]:
    """Helper function for importing income-, age-, education-,
    and occupation-tables in one call.

    dropna: Call dropna() on each table. 
    index_code: Replace generic table index with the area-code.
    """

    income = import_income("income_2022_line.csv")
    age = import_age("age_2022_line.csv")
    education = import_education("education_2022_line.csv")
    occupation = import_occupation("occupation_2022_line.csv")

    if dropna == True:
        income.dropna(inplace=True)
        age.dropna(inplace=True)
        education.dropna(inplace=True)
        occupation.dropna(inplace=True)

    if index_code:
        income.set_index(income['code'], inplace=True)
        age.set_index(age['code'], inplace=True)
        education.set_index(education['code'], inplace=True)
        occupation.set_index(occupation['code'], inplace=True)

    return income, age, education, occupation


if __name__ == "__main__":
    income, age, education, occupation = import_all()

    age_frac = age.drop(['code', 'area', 'male', 'female', 'avg', 'pop'], 
                        axis=1).div(age['pop'], axis=0)

    education_frac = education.drop(['code', 'area', 'over_18'], 
                                    axis=1).div(education['over_18'], axis=0)

    occupation_frac = occupation.drop(['code', 'area', 'pop'], 
                                      axis=1).div(occupation['pop'], axis=0)

    main_table = pd.merge(income, education_frac, 
                         left_index=True, right_index=True)
    main_table = pd.merge(main_table, occupation_frac,
                         left_index=True, right_index=True)

    income_edu_melt = main_table.melt(
        id_vars=['code', 'area', 'avg', 'median', 'low', 'upper'], 
        value_vars=['grade', 'high_school', 'vocational', 'lower_uni', 'higher_uni'], 
        var_name='education_level', 
        value_name='education_frac')
