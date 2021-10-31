#!/bin/bash
assert() {
    expected="$1"
    input="$2"

    ./main "$input" > tmp.s
    cc -o tmp tmp.s
    ./tmp
    actual="$?"

    if [ "$actual" = "$expected" ]; then
        echo "$input => $actual"
    else
        echo "$input => $expected expected, but got $actual"
        exit 1
    fi
}

assert 0 "0;"
assert 42 "42;"
assert 21 "5+20-4;"
assert 41 " 12 + 34 - 5 ;"
assert 47 "5+6*7;"
assert 15 "5*(9-6);"
assert 4 "(3+5)/2;"
assert 10 "-10+20;"
assert 3 "(-1+2)+2;"

assert 1 "0==0;"
assert 1 "1==1;"
assert 0 "0==1;"
assert 0 "1==0;"

assert 1 "1<2;"
assert 0 "1>2;"

assert 1 "1<=2;"
assert 1 "1<=1;"
assert 0 "2<=1;"

assert 1 "2>=1;"
assert 1 "1>=1;"
assert 0 "1>=2;"

assert 1 "1!=2;"
assert 1 "1!=2;"
assert 0 "2!=2;"

assert 1 "a=1; a;"
assert 8 "b=2; b * 4;"
assert 8 "c=2; b=4; c*b;"
assert 9 "return 9;"

assert 1 "abc=1; abc;"
assert 8 "bc=2; bc * 4;"
assert 8 "cd=2; ab=4; cd*ab;"
assert 2 "abc=2; return abc;"

assert 2 "abc=1; if(abc>0) abc=2; abc;"
assert 1 "abc=1; if(abc>1) abc=2; abc;"
assert 2 "abc=1; if(abc>0) abc=2; else abc=3; abc;"
assert 3 "abc=1; if(abc>1) abc=2; else abc=3; abc;"

assert 4 "a=1; while(a<4) a=a+1; a;"
echo OK