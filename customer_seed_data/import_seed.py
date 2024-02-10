import pandas as pd
from pathlib import Path

def import_income(path: str | Path) -> pd.DataFrame:
    """Import income-stats as pandas DataFrame.

    Column-labels:
        'code': Area code
        'area': Area name
        '18':   N Persons over 18
        'avg':  Average income €
        'med':  Median income €
        'low':  N in lowest income bracket
        'mid':  N in middle income bracket
        'upp':  N in upper income bracket
        'cum':  Cumulative purchasing power (?)
    """

    names = ['code', 'area', '18', 'avg', 'med', 'low', 'mid', 'upp', 'cum']
    dtypes = {'code': int,      
              'area': str,      
              '18': 'Int64',    
              'avg': 'Int64',   
              'med': 'Int64',   
              'low': 'Int64',   
              'mid': 'Int64',   
              'upp': 'Int64',   
              'cum': 'Int64'}   

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

    names = ['code', 'area', 'pop', 'male', 'female', 'avg', '18_19', '20_24', 
             '25_29', '30_34', '35_39', '40_44', '45_49', '50_54', '55_59',
             '60_64', '65_69', '70_74', '75_79', '80_84', '85']
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
              '85': 'Int64'}

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
        18:     Population above age of 18.     
        grade:  Finished grade school.
        educated: Finished education higher than grade school.
        high_school: Finished high school.
        vocational: Finished vocational education.
        lower_uni: Finished lower university class education.
        higher_uni: Finished higher university class education.
    """

    names =['code', 'area', '18', 'grade', 'educated', 
            'high_school', 'vocational',
            'lower_uni', 'higher_uni']
    dtypes = {
            'code': int,
            'area': str, 
            '18': 'Int64', 
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


if __name__ == "__main__":
    income = import_income("income_2022_line.csv")
    age = import_age("age_2022_line.csv")
    education = import_education("education_2022_line.csv")
    occupation = import_occupation("occupation_2022_line.csv")