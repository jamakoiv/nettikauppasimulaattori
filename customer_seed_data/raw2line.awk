# Convert CSV-file in format
#
# "12345 Area name", "Data explain text A", 111
# "12345 Area name", "Data explain text B", 222
# "12345 Area name", "Data explain text C", 333
# "12345 Area name", "Data explain text D", 444
# "12345 Area name", "Data explain text E", 555
# "12345 Area name", "Data explain text F", 666
# "12345 Area name", "Data explain text G", 777
#
# to single line per "12345 Area name"
#
# 12345, Area name, 111, 222, 333, 444, 555, 666, 777
#
# Usage: awk -F '<separator>' -v <N> -f <this_file> <input_file>
#
{
    {n++}  # Counter for tracking lines.

    # For first line of the data, output entire line properly formatted.
    if (n == 1) {
        ORS=""
        postal = substr($1, 2, 5)
        area = substr($1, 9, length($1)-9)
        value = substr($3, 0, length($3)-1)  # Remove line-break.

        print postal"; "area"; "value"; "
    } 
    # Last line. Write last entry and newline.
    # N must be supplied from command line with -v
    else if(n == N) {
        ORS="\n"
        print $3;
        n = 0
    } 
    # For lines 2-6 of the data. Write last entry without newline.
    else {  
        ORS=""

        value = substr($3, 0, length($3)-1)
        print value"; "
    }
}
