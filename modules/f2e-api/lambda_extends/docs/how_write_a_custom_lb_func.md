### 如何撰寫一個 Lambda function

1. 設定 `conf/lambdaSetup.json`
  * function_name 名字
  * file_path 位址
  * params 可定義傳遞參數
    * "params_name:type"
  * description 描述
  * ``` {
    "function_name": "sumAll",
    "file_path": "sumAll.js",
    "params": ["aliasName:string"],
    "description": "加總所有的數值"
  }
  ```
2. 在`js/xxx.js` 編寫javascript template
  * `xxx` 和 `function_name`的命名必須一樣
  * 在params 中設定的參數必須給予預設值
    * ex. `aliasName = (typeof aliasName == "undefined"? "SummaryResult" : aliasName)`
  * 資料預設會放在 `input` 這個變數之中
    * 透過 `JSON.stringify(input)` 產生運算的準備資料
  * 最後產生的結果必須放在 `output` 變數之中, 且為array
    * output = JSON.stringify([output_tmp])
3. 在f2e-api/test/api/lambda_query 撰寫測試案例
  * 測試的mock資料在 `data/test_data_sample1.json`
