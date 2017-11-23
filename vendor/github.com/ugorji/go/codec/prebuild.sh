#!/bin/bash

# _needgen is a helper function to tell if we need to generate files for msgp, codecgen.
_needgen() ***REMOVED***
    local a="$1"
    zneedgen=0
    if [[ ! -e "$a" ]]
    then
        zneedgen=1
        echo 1
        return 0
    fi 
    for i in `ls -1 *.go.tmpl gen.go values_test.go`
    do
        if [[ "$a" -ot "$i" ]]
        then
            zneedgen=1
            echo 1
            return 0
        fi 
    done 
    echo 0
***REMOVED***

# _build generates fast-path.go and gen-helper.go.
# 
# It is needed because there is some dependency between the generated code
# and the other classes. Consequently, we have to totally remove the 
# generated files and put stubs in place, before calling "go run" again
# to recreate them.
_build() ***REMOVED***
    if ! [[ "$***REMOVED***zforce***REMOVED***" == "1" ||
                "1" == $( _needgen "fast-path.generated.go" ) ||
                "1" == $( _needgen "gen-helper.generated.go" ) ||
                "1" == $( _needgen "gen.generated.go" ) ||
                1 == 0 ]]
    then
        return 0
    fi 

   # echo "Running prebuild"
    if [ "$***REMOVED***zbak***REMOVED***" == "1" ] 
    then
        # echo "Backing up old generated files"
        _zts=`date '+%m%d%Y_%H%M%S'`
        _gg=".generated.go"
        [ -e "gen-helper$***REMOVED***_gg***REMOVED***" ] && mv gen-helper$***REMOVED***_gg***REMOVED*** gen-helper$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
        [ -e "fast-path$***REMOVED***_gg***REMOVED***" ] && mv fast-path$***REMOVED***_gg***REMOVED*** fast-path$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
        # [ -e "safe$***REMOVED***_gg***REMOVED***" ] && mv safe$***REMOVED***_gg***REMOVED*** safe$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
        # [ -e "unsafe$***REMOVED***_gg***REMOVED***" ] && mv unsafe$***REMOVED***_gg***REMOVED*** unsafe$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
    else 
        rm -f fast-path.generated.go gen.generated.go gen-helper.generated.go \
           *safe.generated.go *_generated_test.go *.generated_ffjson_expose.go
    fi

    cat > gen.generated.go <<EOF
// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// DO NOT EDIT. THIS FILE IS AUTO-GENERATED FROM gen-dec-(map|array).go.tmpl

const genDecMapTmpl = \`
EOF

    cat >> gen.generated.go < gen-dec-map.go.tmpl

    cat >> gen.generated.go <<EOF
\`

const genDecListTmpl = \`
EOF

    cat >> gen.generated.go < gen-dec-array.go.tmpl

    cat >> gen.generated.go <<EOF
\`

EOF

    cat > gen-from-tmpl.codec.generated.go <<EOF
package codec 
import "io"
func GenInternalGoFile(r io.Reader, w io.Writer, safe bool) error ***REMOVED***
return genInternalGoFile(r, w, safe)
***REMOVED***
EOF
    
    cat > gen-from-tmpl.generated.go <<EOF
//+build ignore

package main

//import "flag"
import "ugorji.net/codec"
import "os"

func run(fnameIn, fnameOut string, safe bool) ***REMOVED***
fin, err := os.Open(fnameIn)
if err != nil ***REMOVED*** panic(err) ***REMOVED***
defer fin.Close()
fout, err := os.Create(fnameOut)
if err != nil ***REMOVED*** panic(err) ***REMOVED***
defer fout.Close()
err = codec.GenInternalGoFile(fin, fout, safe)
if err != nil ***REMOVED*** panic(err) ***REMOVED***
***REMOVED***

func main() ***REMOVED***
// do not make safe/unsafe variants. 
// Instead, depend on escape analysis, and place string creation and usage appropriately.
// run("unsafe.go.tmpl", "safe.generated.go", true)
// run("unsafe.go.tmpl", "unsafe.generated.go", false)
run("fast-path.go.tmpl", "fast-path.generated.go", false)
run("gen-helper.go.tmpl", "gen-helper.generated.go", false)
***REMOVED***

EOF
    go run -tags=notfastpath gen-from-tmpl.generated.go && \
        rm -f gen-from-tmpl.*generated.go 
***REMOVED***

_codegenerators() ***REMOVED***
    if [[ $zforce == "1" || 
                "1" == $( _needgen "values_codecgen$***REMOVED***zsfx***REMOVED***" ) ||
                "1" == $( _needgen "values_msgp$***REMOVED***zsfx***REMOVED***" ) ||
                "1" == $( _needgen "values_ffjson$***REMOVED***zsfx***REMOVED***" ) ||
                1 == 0 ]] 
    then
        # codecgen creates some temporary files in the directory (main, pkg).
        # Consequently, we should start msgp and ffjson first, and also put a small time latency before
        # starting codecgen.
        # Without this, ffjson chokes on one of the temporary files from codecgen.
        if [[ $zexternal == "1" ]]
        then 
            echo "ffjson ... " && \
                ffjson -w values_ffjson$***REMOVED***zsfx***REMOVED*** $zfin &
            zzzIdFF=$!
            echo "msgp ... " && \
                msgp -tests=false -o=values_msgp$***REMOVED***zsfx***REMOVED*** -file=$zfin &
            zzzIdMsgp=$!
            
            sleep 1 # give ffjson and msgp some buffer time. see note above.
        fi
        
        echo "codecgen - !unsafe ... " && \
            codecgen -rt codecgen -t 'x,codecgen,!unsafe' -o values_codecgen$***REMOVED***zsfx***REMOVED*** -d 19780 $zfin &
        zzzIdC=$!
        echo "codecgen - unsafe ... " && \
            codecgen  -u -rt codecgen -t 'x,codecgen,unsafe' -o values_codecgen_unsafe$***REMOVED***zsfx***REMOVED*** -d 19781 $zfin &
        zzzIdCU=$!
        wait $zzzIdC $zzzIdCU $zzzIdMsgp $zzzIdFF && \
            # remove (M|Unm)arshalJSON implementations, so they don't conflict with encoding/json bench \
            if [[ $zexternal == "1" ]]
            then
                sed -i 's+ MarshalJSON(+ _MarshalJSON(+g' values_ffjson$***REMOVED***zsfx***REMOVED*** && \
                    sed -i 's+ UnmarshalJSON(+ _UnmarshalJSON(+g' values_ffjson$***REMOVED***zsfx***REMOVED***
            fi && \
            echo "generators done!" && \
            true
    fi 
***REMOVED***

# _init reads the arguments and sets up the flags
_init() ***REMOVED***
OPTIND=1
while getopts "fbx" flag
do
    case "x$flag" in 
        'xf') zforce=1;;
        'xb') zbak=1;;
        'xx') zexternal=1;;
        *) echo "prebuild.sh accepts [-fbx] only"; return 1;;
    esac
done
shift $((OPTIND-1))
OPTIND=1
***REMOVED***

# main script.
# First ensure that this is being run from the basedir (i.e. dirname of script is .)
if [ "." = `dirname $0` ]
then
    zmydir=`pwd`
    zfin="test_values.generated.go"
    zsfx="_generated_test.go"
    # rm -f *_generated_test.go 
    rm -f codecgen-*.go && \
        _init "$@" && \
        _build && \
        cp $zmydir/values_test.go $zmydir/$zfin && \
        _codegenerators && \
        echo prebuild done successfully
    rm -f $zmydir/$zfin
else
    echo "Script must be run from the directory it resides in"
fi 

