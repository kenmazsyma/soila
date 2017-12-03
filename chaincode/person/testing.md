# Testing

## Test case

|No|target|pattern|test function|
|---|:--|:--|:--|
|a-1|Register|success|Test_Register|
|a-2||number of parameters<expect|Test_Register|
|a-3||number of parameters>expect|Test_Register|
|a-4||duplicate|Test_Register|
|b-1|Update|success|Test_Register|
|b-2||number of parameters<expect|Test_Register|
|b-3||number of parameters>expect|Test_Register|
|b-4||not found|Test_Register|
|b-5||not owned|Test_Register|
|b-6||same data|Test_Register|
|c-1|Get|success|Test_Register|
|d-1|AddActivity|success|Test_Register|
|d-2||number of parameters<expect|Test_Register|
|d-3||number of parameters>expect|Test_Register|
|d-4||not found|Test_Register|
|d-5||not owned|Test_Register|
|e-1|AddReputation|success|Test_Register|
|e-2||number of parameters<expect|Test_Register|
|e-3||number of parameters>expect|Test_Register|
|e-4||not found|Test_Register|
|f-1|RemoveReputation|success|not tested|
|f-2||number of parameters<expect|not tested|
|f-3||number of parameters>expect|not tested|
|f-4||not found|not tested|
