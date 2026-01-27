#!/usr/bin/env bash

# Получаем аргументы командной строки
cpu_profile=${1:-"base_cpu.pprof"}
mem_profile=${2:-"base_mem.pprof"}

# Запускаем тесты Go с профилированием CPU и памяти
go test -bench=. \
          -benchmem \
          -cpuprofile="${cpu_profile}" \
          -memprofile="${mem_profile}"