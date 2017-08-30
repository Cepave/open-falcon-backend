if [ "$1" == "init" ]; then
  echo "strategy_cgt_test.go"
  go test -v strategy_cgt_test.go -test.run TestStrategyCreate
fi
echo "strategy_cgt_test.go"
go test -v strategy_cgt_test.go -test.run TestStrategyGet
