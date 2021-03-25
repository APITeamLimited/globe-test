#!/bin/bash

# Run all the different permutations of all the tests and other things
# This helps ensure that nothing gets broken.

_tests() ***REMOVED***
    local vet="" # TODO: make it off
    local gover=$( go version | cut -f 3 -d ' ' )
    # note that codecgen requires fastpath, so you cannot do "codecgen notfastpath"
    local a=( "" "safe"  "notfastpath" "notfastpath safe" "codecgen" "codecgen safe" )
    for i in "$***REMOVED***a[@]***REMOVED***"
    do
        echo ">>>> TAGS: $i"
        local i2=$***REMOVED***i:-default***REMOVED***
        case $gover in
            go1.[0-6]*) go test $***REMOVED***zargs[*]***REMOVED*** -tags "$i" "$@" ;;
            *) go vet -printfuncs "errorf" "$@" &&
                     go test $***REMOVED***zargs[*]***REMOVED*** -vet "$vet" -tags "alltests $i" -run "Suite" -coverprofile "$***REMOVED***i2// /-***REMOVED***.cov.out" "$@" ;;
        esac
        if [[ "$?" != 0 ]]; then return 1; fi 
    done
    echo "++++++++ TEST SUITES ALL PASSED ++++++++"
***REMOVED***


# is a generation needed?
_ng() ***REMOVED***
    local a="$1"
    if [[ ! -e "$a" ]]; then echo 1; return; fi 
    for i in `ls -1 *.go.tmpl gen.go values_test.go`
    do
        if [[ "$a" -ot "$i" ]]; then echo 1; return; fi 
    done
***REMOVED***

_prependbt() ***REMOVED***
    cat > $***REMOVED***2***REMOVED*** <<EOF
// +build generated

EOF
    cat $***REMOVED***1***REMOVED*** >> $***REMOVED***2***REMOVED***
    rm -f $***REMOVED***1***REMOVED***
***REMOVED***

# _build generates fast-path.go and gen-helper.go.
_build() ***REMOVED***
    if ! [[ "$***REMOVED***zforce***REMOVED***" || $(_ng "fast-path.generated.go") || $(_ng "gen-helper.generated.go") || $(_ng "gen.generated.go") ]]; then return 0; fi 
    
    if [ "$***REMOVED***zbak***REMOVED***" ]; then
        _zts=`date '+%m%d%Y_%H%M%S'`
        _gg=".generated.go"
        [ -e "gen-helper$***REMOVED***_gg***REMOVED***" ] && mv gen-helper$***REMOVED***_gg***REMOVED*** gen-helper$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
        [ -e "fast-path$***REMOVED***_gg***REMOVED***" ] && mv fast-path$***REMOVED***_gg***REMOVED*** fast-path$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
        [ -e "gen$***REMOVED***_gg***REMOVED***" ] && mv gen$***REMOVED***_gg***REMOVED*** gen$***REMOVED***_gg***REMOVED***__$***REMOVED***_zts***REMOVED***.bak
    fi 
    rm -f gen-helper.generated.go fast-path.generated.go gen.generated.go \
       *safe.generated.go *_generated_test.go *.generated_ffjson_expose.go 

    cat > gen.generated.go <<EOF
// +build codecgen.exec

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
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

const genEncChanTmpl = \`
EOF
    cat >> gen.generated.go < gen-enc-chan.go.tmpl
    cat >> gen.generated.go <<EOF
\`
EOF
    cat > gen-from-tmpl.codec.generated.go <<EOF
package codec 
import "io"
func GenInternalGoFile(r io.Reader, w io.Writer) error ***REMOVED***
return genInternalGoFile(r, w)
***REMOVED***
EOF
    cat > gen-from-tmpl.generated.go <<EOF
//+build ignore

package main

import "$***REMOVED***zpkg***REMOVED***"
import "os"

func run(fnameIn, fnameOut string) ***REMOVED***
println("____ " + fnameIn + " --> " + fnameOut + " ______")
fin, err := os.Open(fnameIn)
if err != nil ***REMOVED*** panic(err) ***REMOVED***
defer fin.Close()
fout, err := os.Create(fnameOut)
if err != nil ***REMOVED*** panic(err) ***REMOVED***
defer fout.Close()
err = codec.GenInternalGoFile(fin, fout)
if err != nil ***REMOVED*** panic(err) ***REMOVED***
***REMOVED***

func main() ***REMOVED***
run("fast-path.go.tmpl", "fast-path.generated.go")
run("gen-helper.go.tmpl", "gen-helper.generated.go")
run("mammoth-test.go.tmpl", "mammoth_generated_test.go")
run("mammoth2-test.go.tmpl", "mammoth2_generated_test.go")
// run("sort-slice.go.tmpl", "sort-slice.generated.go")
***REMOVED***
EOF

    sed -e 's+// __DO_NOT_REMOVE__NEEDED_FOR_REPLACING__IMPORT_PATH__FOR_CODEC_BENCH__+import . "github.com/ugorji/go/codec"+' \
        shared_test.go > bench/shared_test.go

    # explicitly return 0 if this passes, else return 1
    go run -tags "prebuild" prebuild.go || return 1
    go run -tags "notfastpath safe codecgen.exec" gen-from-tmpl.generated.go || return 1
    rm -f gen-from-tmpl.*generated.go
    return 0
***REMOVED***

_codegenerators() ***REMOVED***
    local c5="_generated_test.go"
    local c7="$PWD/codecgen"
    local c8="$c7/__codecgen"
    local c9="codecgen-scratch.go"

    if ! [[ $zforce || $(_ng "values_codecgen$***REMOVED***c5***REMOVED***") ]]; then return 0; fi
    
    # Note: ensure you run the codecgen for this codebase/directory i.e. ./codecgen/codecgen
    true &&
        echo "codecgen ... " &&
        if [[ $zforce || ! -f "$c8" || "$c7/gen.go" -nt "$c8" ]]; then
            echo "rebuilding codecgen ... " && ( cd codecgen && go build -o $c8 $***REMOVED***zargs[*]***REMOVED*** . )
        fi &&
        $c8 -rt codecgen -t 'codecgen generated' -o values_codecgen$***REMOVED***c5***REMOVED*** -d 19780 $zfin $zfin2 &&
        cp mammoth2_generated_test.go $c9 &&
        $c8 -t 'codecgen,!notfastpath generated,!notfastpath' -o mammoth2_codecgen$***REMOVED***c5***REMOVED*** -d 19781 mammoth2_generated_test.go &&
        rm -f $c9 &&
        echo "generators done!" 
***REMOVED***

_prebuild() ***REMOVED***
    echo "prebuild: zforce: $zforce"
    local d="$PWD"
    local zfin="test_values.generated.go"
    local zfin2="test_values_flex.generated.go"
    local zpkg="github.com/ugorji/go/codec"
    # zpkg=$***REMOVED***d##*/src/***REMOVED***
    # zgobase=$***REMOVED***d%%/src/****REMOVED***
    # rm -f *_generated_test.go 
    rm -f codecgen-*.go &&
        _build &&
        cp $d/values_test.go $d/$zfin &&
        cp $d/values_flex_test.go $d/$zfin2 &&
        _codegenerators &&
        if [[ "$(type -t _codegenerators_external )" = "function" ]]; then _codegenerators_external ; fi &&
        if [[ $zforce ]]; then go install $***REMOVED***zargs[*]***REMOVED*** .; fi &&
        echo "prebuild done successfully"
    rm -f $d/$zfin $d/$zfin2
    # unset zfin zfin2 zpkg
***REMOVED***

_make() ***REMOVED***
    local makeforce=$***REMOVED***zforce***REMOVED***
    zforce=1
    (cd codecgen && go install $***REMOVED***zargs[*]***REMOVED*** .) && _prebuild && go install $***REMOVED***zargs[*]***REMOVED*** .
    zforce=$***REMOVED***makeforce***REMOVED***
***REMOVED***

_clean() ***REMOVED***
    rm -f gen-from-tmpl.*generated.go \
       codecgen-*.go \
       test_values.generated.go test_values_flex.generated.go
***REMOVED***

_release() ***REMOVED***
    local reply
    read -p "Pre-release validation takes a few minutes and MUST be run from within GOPATH/src. Confirm y/n? " -n 1 -r reply
    echo
    if [[ ! $reply =~ ^[Yy]$ ]]; then return 1; fi

    # expects GOROOT, GOROOT_BOOTSTRAP to have been set.
    if [[ -z "$***REMOVED***GOROOT// ***REMOVED***" || -z "$***REMOVED***GOROOT_BOOTSTRAP// ***REMOVED***" ]]; then return 1; fi
    # (cd $GOROOT && git checkout -f master && git pull && git reset --hard)
    (cd $GOROOT && git pull)
    local f=`pwd`/make.release.out
    cat > $f <<EOF
========== `date` ===========
EOF
    # # go 1.6 and below kept giving memory errors on Mac OS X during SDK build or go run execution,
    # # that is fine, as we only explicitly test the last 3 releases and tip (2 years).
    local makeforce=$***REMOVED***zforce***REMOVED***
    zforce=1
    for i in 1.10 1.11 1.12 master
    do
        echo "*********** $i ***********" >>$f
        if [[ "$i" != "master" ]]; then i="release-branch.go$i"; fi
        (false ||
             (echo "===== BUILDING GO SDK for branch: $i ... =====" &&
                  cd $GOROOT &&
                  git checkout -f $i && git reset --hard && git clean -f . &&
                  cd src && ./make.bash >>$f 2>&1 && sleep 1 ) ) &&
            echo "===== GO SDK BUILD DONE =====" &&
            _prebuild &&
            echo "===== PREBUILD DONE with exit: $? =====" &&
            _tests "$@"
        if [[ "$?" != 0 ]]; then return 1; fi
    done
    zforce=$***REMOVED***makeforce***REMOVED***
    echo "++++++++ RELEASE TEST SUITES ALL PASSED ++++++++"
***REMOVED***

_usage() ***REMOVED***
    cat <<EOF
primary usage: $0 
    -[tmpfxnld]           -> [tests, make, prebuild (force) (external), inlining diagnostics, mid-stack inlining, race detector]
    -v                    -> verbose
EOF
    if [[ "$(type -t _usage_run)" = "function" ]]; then _usage_run ; fi
***REMOVED***

_main() ***REMOVED***
    if [[ -z "$1" ]]; then _usage; return 1; fi
    local x
    local zforce
    local zargs=()
    local zverbose=()
    local zbenchflags=""
    OPTIND=1
    while getopts ":ctmnrgpfvlyzdb:" flag
    do
        case "x$flag" in
            'xf') zforce=1 ;;
            'xv') zverbose+=(1) ;;
            'xl') zargs+=("-gcflags"); zargs+=("-l=4") ;;
            'xn') zargs+=("-gcflags"); zargs+=("-m=2") ;;
            'xd') zargs+=("-race") ;;
            'xb') x='b'; zbenchflags=$***REMOVED***OPTARG***REMOVED*** ;;
            x\?) _usage; return 1 ;;
            *) x=$flag ;;
        esac
    done
    shift $((OPTIND-1))
    # echo ">>>> _main: extra args: $@"
    case "x$x" in
        'xt') _tests "$@" ;;
        'xm') _make "$@" ;;
        'xr') _release "$@" ;;
        'xg') _go ;;
        'xp') _prebuild "$@" ;;
        'xc') _clean "$@" ;;
        'xy') _analyze_extra "$@" ;;
        'xz') _analyze "$@" ;;
        'xb') _bench "$@" ;;
    esac
    # unset zforce zargs zbenchflags
***REMOVED***

[ "." = `dirname $0` ] && _main "$@"

