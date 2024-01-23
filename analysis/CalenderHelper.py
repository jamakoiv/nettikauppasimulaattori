#!/usr/bin/python3

from datetime import timedelta, datetime
from typing import Tuple

import pandas as pd
import numpy as np
import matplotlib as mpl
import matplotlib.dates


def GetWeekEndpoints(date: datetime,
                     end_offset: int = 0) -> Tuple[datetime, datetime]:
    """
    Given a date, return datetime-objects corresponding to the
    start and end of the week containing the date.

    date: datetime- or date-object.
    end_offset: 0 if you want monday to sunday,
                1 if you want monday to monday.
    """

    start = date - timedelta(days=date.weekday())
    end = date + timedelta(days=end_offset+6-date.weekday())

    return start, end


def GetHistogramBins(date_start: datetime,
                     date_end: datetime,
                     freq: str,
                     labels_format: str = "%Y-%m-%d") -> Tuple[np.ndarray, list[str]]:
    """
    Return histogram bins suitable for usage with matplotlib histogram plot.

    Input:
    --------
    date_start: datetime-object for start date.
    date_end: datetime-object for end date.
    freq: pandas time offset alias, eq. '1D' for one day, '1H' for on hour etc.

    Output:
    --------
    numpy-array containing bins.
    list containing labels for the bins.
    """

    tmp = pd.date_range(date_start, date_end, freq=freq)
    bins = mpl.dates.date2num(tmp)
    labels = tmp.strftime(labels_format)

    return bins, labels


if __name__ == "__main__":
    pass
