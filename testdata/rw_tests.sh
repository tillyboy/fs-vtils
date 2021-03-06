#!/usr/bin/bash

tests_rw=$src/testdata/rw

function reset_rw {
  cd $tests_rw

  for i in {1..11}; do

    if [[ $i -lt 10 ]]; then
      c=case0$i
    else
      c=case$i
    fi

    rm -R $c
    mkdir $c

    echo -e "a\nb\nc" > $c/file
    touch $c/empty

  done

  cd $src
}

function evaluate_rw {
  cd $tests_rw

  #case01
  c=case01
  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="a\nb\nc" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ ! $(cat $c/empty)="written" ]]; then
    fail "rw" $c 4
  fi

  #case02
  c=case02
  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="written" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ ! $(cat $c/empty)="" ]]; then
    fail "rw" $c 4
  fi

  #case03
  c=case03
  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="a\nb\nc" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ -ne $(cat $c/empty) ]]; then
    fail "rw" $c 4
  # elif [[ ! -f $c/new ]]; then
  #   fail "rw" $c 5
  # elif [[ ! $(cat $c/new)="written" ]]; then
  #   fail "rw" $c 6
  fi


  #case04
  unchanged_rw case04

  #case05
  unchanged_rw case05

  #case06
  unchanged_rw case06

  c=case07
  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="a\nb\nc" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ ! $(cat $c/empty)="appended" ]]; then
    fail "rw" $c 4
  fi

  c=case08
  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="a\nb\nc\nappended" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ ! $(cat $c/empty)="" ]]; then
    fail "rw" $c 4
  fi

  unchanged_rw case09

  c=case10
  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="a\nb\nc" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ ! $(cat $c/empty)="a\nb\nc" ]]; then
    fail "rw" $c 4
  fi

  unchanged_rw case11

  cd $src
}

function unchanged_rw {
  local c=$1

  if [[ ! -f $c/file ]]; then
    fail "rw" $c 1
  elif [[ ! $(cat $c/file)="a\nb\nc" ]]; then
    fail "rw" $c 2
  elif [[ ! -f $c/empty ]]; then
    fail "rw" $c 3
  elif [[ -ne $(cat $c/empty) ]]; then
    fail "rw" $c 4
  fi
}
