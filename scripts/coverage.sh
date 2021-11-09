#!/bin/bash
echo "Quality Gate: checking if test coverage is above threshold ..."
echo "Threshold             : ${TESTCOVERAGE_THRESHOLD} %"
totalCoverage=`go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'`
echo "Current test coverage : $totalCoverage %"
if (( $(echo "$totalCoverage ${TESTCOVERAGE_THRESHOLD}" | awk '{print ($1 >= $2)}') )); then
    echo "OK"
else
    echo "Current test coverage is below threshold. Please add more tests"
    echo "Failed"
    exit 1
fi