aliasName = (typeof aliasName == "undefined"? "SummaryResult" : aliasName)
output = JSON.stringify(input)
//for solve null value of open-falcon, generate a clean object
var output_tmp = {}
output_tmp.Counter = input[0].Counter
output_tmp.DsType = input[0].DsType
output_tmp.Endpoint = aliasName
output_tmp.Step = input[0].Step
var values = []
var tmpR = _.filter(input, function(obj){
  return obj.Values.length != 0
})

if(tmpR.length != 0){
  _.each(tmpR[0].Values, function(val, indx){
    value = { Value: 0,
              Timestamp: val.Timestamp }
    values.push(value)
    return
  })
  output_tmp.Values = values
  //compute
  _.each(input, function(record){
    _.each(record.Values, function(v, indx){
      value = (isNaN(v.Value)? 0 : v.Value)
      return output_tmp.Values[indx].Value += value
    })
    return record
  })
  output = JSON.stringify([output_tmp])
}
