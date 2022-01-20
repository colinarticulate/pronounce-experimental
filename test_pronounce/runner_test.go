package main

import (
  "fmt"
  // "io/ioutil"
  "reflect"
  "testing"
)

func Test_parseSimpleExpectation(te *testing.T) {
  // GIVEN
  tokens := []string{
    "t", "th", "good", "possible",
  }

  // WHEN 
  actual, remaining, err := parseSimpleExpectation(tokens)
  expected := simpleExpect{
    []string{
      "t", "th",
    },
    []string{
      "good", "possible",
    },
  }
  
  // THEN 
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{}) {
    te.Error()
  }
  
  // and GIVEN 
  tokens = []string{
    "t", "th", "good", "possible", "(",
  }
  
  // WHEN 
  actual, remaining, err = parseSimpleExpectation(tokens)
  expected = simpleExpect{
    []string{
      "t", "th",
    },
    []string{
      "good", "possible",
    },
  }
  
  // THEN 
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{"("}) {
    te.Error()
  }
  
  // and GIVEN 
  tokens = []string{
    "surprise", "t", "th", "good", "possible", "(",
  }
  
  // WHEN 
  _, _, err = parseSimpleExpectation(tokens)

  // THEN 
  if err != parseError {
    te.Error()
  }
  
  // and GIVEN a malformed set of tokens
  tokens = []string{
    "surprise", "t", "(", "th", "good", "possible", "(",
  }
  
  // WHEN 
  _, _, err = parseSimpleExpectation(tokens)

  // THEN 
  if err != parseError {
    te.Error()
  }
}

func Test_parseOrExpectation(te *testing.T) {
  // GIVEN 
  tokens := []string{
    "(", "t", "th", "good", "possible", "||", "eh", "good", ")",
  }
  
  // WHEN 
  actual, remaining, err := parseOrExpectation(tokens)
  expected := orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "t", "th",
          },
          []string{
            "good", "possible",
          },        
        },
      },
      {
        {
          []string{
            "eh",
          },
          []string{
            "good",
          },      
        },        
      }, 
    },
  }
  
  // THEN 
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{}) {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN 
  tokens = []string{
    "(", "t", "good", "th", "possible", "||", "eh", "good", ")",
  }
  
  // WHEN 
  actual, remaining, err = parseOrExpectation(tokens)
  expected = orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "t",
          },
          []string{
            "good",
          },        
        },
        {
          []string{
            "th",
          },
          []string{
            "possible",
          },        
        },
      },
      {
        {
          []string{
            "eh",
          },
          []string{
            "good",
          },      
        },        
      }, 
    },
  }
  
  // THEN
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{}) {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN 
  tokens = []string{
    "(", "t", "good", "th", "possible", "||", "eh", "good", "||", "iy", "possible", ")",
  }
  
  // WHEN 
  actual, remaining, err = parseOrExpectation(tokens)
  expected = orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "t",
          },
          []string{
            "good",
          },        
        },
        {
          []string{
            "th",
          },
          []string{
            "possible",
          },        
        },
      },
      {
        {
          []string{
            "eh",
          },
          []string{
            "good",
          },      
        },        
      }, 
      {
        {
          []string{
            "iy",
          },
          []string{
            "possible",
          },      
        },        
      }, 
    },
  }
  
  // THEN
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{}) {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN an unterminated or of tokens
  tokens = []string{
    "(", "t", "good", "th", "possible", "||", "eh", "good",
  }
  
  // WHEN 
  _, _, err = parseOrExpectation(tokens)

  // THEN 
  if err == nil {
    te.Error()
  }
  
  // and GIVEN 
  tokens = []string{
    "(", "t", "good", "th", "possible", "||", "eh", "good", ")", "junk",
  }
  
  // WHEN 
  _, remaining, err = parseOrExpectation(tokens)

  // THEN
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{"junk"}) {
    te.Error()
  }
}

func Test_notExpectation(te *testing.T) {
  // GIVEN 
  tokens := []string{
    "~", "t", "th", "good", "possible",
  }
  
  // WHEN 
  actual, remaining, err := parseNotExpectation(tokens)
  expected := notExpect{
    simpleExpect{
      []string{
        "t", "th",
      },
      []string{
        "good", "possible",
      },       
    },  
  }
  
  // THEN 
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{}) {
    te.Error()
  }
  
  // and GIVEN
  tokens = []string{
    "~", "t", "th", "good", "possible", "eh", "good",
  }
  
  // WHEN 
  actual, remaining, err = parseNotExpectation(tokens)
  expected = notExpect{
    simpleExpect{
      []string{
        "t", "th",
      },
      []string{
        "good", "possible",
      },       
    },  
  }
  
  // THEN 
  if err != nil {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  if !reflect.DeepEqual(remaining, []string{"eh", "good",}) {
    te.Error()
  }
}

func Test_fetchExpectations(te *testing.T) {
  // GIVEN 
  file := "test/expectations.txt"
  
  // WHEN 
  actual, err := fetchExpectations(file)
  expected := testsOut{
    {
      "file1",
      []result{},
      []expectation{
        simpleExpect{
          []string{
            "t", "th",
          },
          []string{
            "good", "possible",
          },                  
        },
        notExpect{
          simpleExpect{
            []string{
              "eh",
            },
            []string{
              "absent",
            },                  
          },          
        },
        simpleExpect{
          []string{
            "eh",
          },
          []string{
            "missing",
          },        
        },
        orExpect{
          [][]simpleExpect{
            {
              {
                []string{
                  "*",
                },
                []string{
                  "surprise",
                },        
              },
              {
                []string{
                  "s",
                },
                []string{
                  "good",
                },        
              },
            },
            {
              {
                []string{
                  "s",
                },
                []string{
                  "missing",
                },      
              },        
            }, 
          },          
        },
        simpleExpect{
          []string{
            "t",
          },
          []string{
            "missing",
          },        
        },
      },
      false,      
    },
  }
  
  // THEN 
  if err != nil {
    te.Error(err)
  }
  if !reflect.DeepEqual(actual, expected) {
    fmt.Println(actual)
    te.Error()
  }
}

func Test_parseResult(t *testing.T) {
  // GIVEN
  tokens := []string{
    "t", "good", "eh", "missing",  "s", "possible", "t", "good",
  }
  
  // WHEN 
  actual := parseResult(tokens)
  expected := []result{
    {
      "t", 
      "good", 
    },
    {
      "eh", 
      "missing", 
    },
    {
      "s", 
      "possible", 
    },
    {
      "t", 
      "good", 
    },
  }
  
  // THEN 
  if !reflect.DeepEqual(actual, expected) {
    t.Error()
  }
}

func Test_parseTestExpectation(t *testing.T) {
  // GIVEN
  tokens := []string{
    "(", "*", "surprise", "*", "surprise", "absent", "s", "good", "possible", "missing", "||", "s", "missing", "*", "surprise", "absent", "*", "surprise", "absent", "||", "*", "surprise", "absent", "s", "missing", "*", "surprise", "absent", ")", "k", "good", "possible", "uw", "good", "possible", "l", "good", "possible",                 
  }
  
  // WHEN 
  actual, err := parseTestExpectation(tokens)
  expected := []expectation{
    orExpect{
      [][]simpleExpect{
        {
          {
            []string{
              "*",
            },
            []string{
              "surprise",
            },
          },
          {
            []string{
              "*",
            },
            []string{
              "surprise", "absent",
            },
          },
          {
            []string{
              "s",
            },
            []string{
              "good", "possible", "missing",
            },
          },
        },
        {
          {
            []string{
              "s",
            },
            []string{
              "missing",
            },
          },
          {
            []string{
              "*",
            },
            []string{
              "surprise", "absent",
            },
          },
          {
            []string{
              "*",
            },
            []string{
              "surprise", "absent",
            },
          },              
        },
        {
          {
            []string{
              "*",
            },
            []string{
              "surprise", "absent",
            },
          },              
          {
            []string{
              "s",
            },
            []string{
              "missing",
            },
          },
          {
            []string{
              "*",
            },
            []string{
              "surprise", "absent",
            },
          },              
        },
      },
    },
    simpleExpect{
      []string{
        "k",
      },
      []string{
        "good", "possible",
      },                  
    },
    simpleExpect{
      []string{
        "uw",
      },
      []string{
        "good", "possible",
      },                  
    },
    simpleExpect{
      []string{
        "l",
      },
      []string{
        "good", "possible",
      },                  
    },
  }
  
  // THEN 
  if err != nil {
    t.Error(err)
  } else if !reflect.DeepEqual(actual, expected) {
    t.Error()
  }
}

func Test_pass(te *testing.T) {
  // Simple expectations...
  // GIVEN 
  result := results{
    {
      "t",
      "possible",
    },
    {
      "eh",
      "missing",
    },
    {
      "sh",
      "surprise",
    },
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },
  }
  expect := simpleExpect{
    []string{
      "t",
    },
    []string{
      "good", "possible",
    },
  }
  
  // WHEN 
  actual, passed := expect.pass(result)
  expected := results{
    {
      "eh",
      "missing",
    },
    {
      "sh",
      "surprise",
    },
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },
  }
  
  // THEN
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN 
  expect = simpleExpect{
    []string{
      "*",
    },
    []string{
      "surprise",
    },
  }
  
  // WHEN 
  actual, passed = expect.pass(result[2:])
  expected = results{
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },
  }
  
  // THEN 
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }  
  
  // and GIVEN
  expect = simpleExpect{
    []string{
      "th",
    },
    []string{
      "absent",
    },
  }
  
  // WHEN 
  actual, passed = expect.pass(result[5:])
  expected = results{
  }

  // THEN
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // Or expectations...
  // GIVEN 
  expectOr := orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "*",
          },
          []string{
            "surprise",
          },        
        },
        {
          []string{
            "s",
          },
          []string{
            "missing",
          },        
        },
      },
      {
        {
          []string{
            "s",
          },
          []string{
            "good",
          },      
        },        
      }, 
    },  
  }
  
  // WHEN 
  actual, passed = expectOr.pass(result[2:])
  expected = results{
    {
      "t",
      "good",
    },
  }
  
  // THEN
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN 
  expectOr = orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "*",
          },
          []string{
            "surprise",
          },        
        },
        {
          []string{
            "s",
          },
          []string{
            "good",
          },        
        },
      },
      {
        {
          []string{
            "sh",
          },
          []string{
            "surprise",
          },      
        },        
      }, 
    },  
  }
  
  // WHEN 
  actual, passed = expectOr.pass(result[2:])
  expected= results{
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },
  }
  
  // THEN 
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN
  expectOr = orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "*",
          },
          []string{
            "surprise",
          },        
        },
        {
          []string{
            "s",
          },
          []string{
            "good",
          },        
        },
      },
      {
        {
          []string{
            "sh",
          },
          []string{
            "possible",
          },      
        },        
      }, 
    },  
  }
  
  // WHEN 
  actual, passed = expectOr.pass(result[2:])
  expected= results{
    {
      "sh",
      "surprise",
    },
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },
  }
  
  // THEN 
  if passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // Not expectations...
  // GIVEN 
  expectNot := notExpect{
    simpleExpect{
      []string{
        "th",
      },
      []string{
        "possible",
      },       
    },      
  }
  
  // WHEN 
  actual, passed = expectNot.pass(result)
  expected = results{
    {
      "eh",
      "missing",
    },
    {
      "sh",
      "surprise",
    },
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },    
  }
  
  // THEN 
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
  
  // and GIVEN 
  expectNot = notExpect{
    simpleExpect{
      []string{
        "t",
      },
      []string{
        "good",
      },       
    },      
  }
  
  // WHEN 
  actual, passed = expectNot.pass(result)
  expected = results{
    {
      "eh",
      "missing",
    },
    {
      "sh",
      "surprise",
    },
    {
      "s",
      "missing",
    },
    {
      "t",
      "good",
    },    
  }
  
  // THEN 
  if !passed {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
}

func  Test_checkResult(t *testing.T) {
  // GIVEN 
  tests := []testOut{
    {
      "audiofile",
      []result{
        {
          "t", 
          "good", 
        },
        {
          "eh", 
          "missing", 
        },
        {
          "s", 
          "possible", 
        },
        {
          "t", 
          "good", 
        },        
      },
      []expectation{
        simpleExpect{
          []string{
            "t",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "ey",
          },
          []string{
            "absent",
          },        
        },
        simpleExpect{
          []string{
            "eh",
          },
          []string{
            "missing",
          },        
        },
        simpleExpect{
          []string{
            "s",
          },
          []string{
            "possible",
          },        
        },
        simpleExpect{
          []string{
            "t",
          },
          []string{
            "good",
          },        
        },
      },
      false,
    },
  }
  
  // WHEN 
  actual := tests[0].checkResult()
  
  // THEN 
  if actual != true {
    t.Error()
  }
  
  // and GIVEN 
  tests = []testOut{
    {
      "audiofile",
      []result{
        {
          "ae", 
          "good", 
        },
        {
          "n", 
          "possible", 
        },
        {
          "ih", 
          "missing", 
        },
        {
          "m", 
          "possible", 
        },        
        {
          "ah", 
          "missing", 
        },        
        {
          "ow", 
          "surprise", 
        },        
        {
          "l", 
          "good", 
        },        
      },
      []expectation{
        simpleExpect{
          []string{
            "ae",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "n",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "ih",
          },
          []string{
            "good",
          },        
        },
        simpleExpect{
          []string{
            "m",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "aa", "aa",
          },
          []string{
            "good", "possible", "absent",
          },        
        },
        simpleExpect{
          []string{
            "l",
          },
          []string{
            "good", "possible",
          },        
        },
      },
      false,
    },    
  }

  // WHEN 
  actual = tests[0].checkResult()

  // THEN 
  if actual != false {
    t.Error()
  }
  
  // and GIVEN 
  tests = []testOut{
    {
      "audiofile",
      []result{
        {
          "hh", 
          "good", 
        },
        {
          "aa", 
          "good", 
        },
        {
          "r", 
          "good", 
        },
        {
          "d", 
          "possible",
        },        
        {
          "l", 
          "missing",
        },      
        {
          "hh", 
          "surprise",
        },       
        {
          "n", 
          "surprise", 
        },   
        {
          "iy", 
          "good",
        },   
      },
      []expectation{
        simpleExpect{
          []string{
            "hh",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "aa",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "r",
          },
          []string{
            "good", "possible", "absent",
          },        
        },
        simpleExpect{
          []string{
            "d",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },        
        },
        simpleExpect{
          []string{
            "l",
          },
          []string{
            "missing",
          },        
        },
        simpleExpect{
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },        
        },
        simpleExpect{
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },        
        },
        simpleExpect{
          []string{
            "iy",
          },
          []string{
            "good", "possible",
          },        
        },
      },
      false,
    },    
  }
  
  // WHEN 
  actual = tests[0].checkResult()
  
  // THEN 
  if actual != true {
    t.Error()
  }
  
  // and GIVEN 
  
  // WHEN 
  actual = tests[0].checkResult()
  
  // THEN 
  if actual != true {
    t.Error()
  }
  
  // and GIVEN 
  tests = []testOut{
    {
      "audiofile",
      []result{
        {
          "ch", 
          "good", 
        },
        {
          "eh", 
          "missing", 
        },
      },
      []expectation{
        simpleExpect{
          []string{
            "ch",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "eh",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "r",
          },
          []string{
            "good", "possible", "absent",
          },        
        },
      },
      false,
    },    
  }
  
  // WHEN 
  actual = tests[0].checkResult()
  
  // THEN 
  if actual == true {
    t.Error()
  }  
  
  // and GIVEN 
  tests = []testOut{
    {
      "audiofile",
      []result{
        {
          "r", 
          "good", 
        },
        {
          "ah", 
          "possible", 
        },
        {
          "sh", 
          "good", 
        },
        {
          "ah", 
          "possible",
        },        
      },
      []expectation{
        simpleExpect{
          []string{
            "r",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "ah",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "sh",
          },
          []string{
            "good", "possible",
          },        
        },
        orExpect{
          [][]simpleExpect{
            {
              {
                []string{
                  "ah", "er",
                },
                []string{
                  "missing",
                },        
              },
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },      
              },        
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },      
              },        
            },
            {
              {
                []string{
                  "ah", "er",
                },
                []string{
                  "good", "possible",
                },        
              },
              {
                []string{
                  "*",
                },
                []string{
                  "surprise",
                },      
              },        
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },      
              },                      
            },
          },
        },
      },
      false,
    },        
  }
  
  // WHEN 
  actual = tests[0].checkResult()
  
  // THEN 
  if actual == true {
    t.Error()
  }    
  
  // and GIVEN
  tests = []testOut{
    {
      "audiofile",
      []result{
        {
          "ae", "surprise",
        },
        {
          "s", "missing",
        },
        {
          "z", "surprise",
        },
        {
          "k", "good",
        },
        {
          "uw", "good",
        },
        {
          "l", "good",
        },
      },
      []expectation{
        orExpect{
          [][]simpleExpect{
            {
              {
                []string{
                  "*",
                },
                []string{
                  "surprise",
                },
              },
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },
              },
              {
                []string{
                  "s",
                },
                []string{
                  "good", "possible", "missing",
                },
              },
            },
            {
              {
                []string{
                  "s",
                },
                []string{
                  "missing",
                },
              },
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },
              },
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },
              },              
            },
            {
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },
              },              
              {
                []string{
                  "s",
                },
                []string{
                  "missing",
                },
              },
              {
                []string{
                  "*",
                },
                []string{
                  "surprise", "absent",
                },
              },              
            },
          },
        },
        simpleExpect{
          []string{
            "k",
          },
          []string{
            "good", "possible",
          },                  
        },
        simpleExpect{
          []string{
            "uw",
          },
          []string{
            "good", "possible",
          },                  
        },
        simpleExpect{
          []string{
            "l",
          },
          []string{
            "good", "possible",
          },                  
        },
      },
      true,
    },
  }
  
  // WHEN 
  actual = tests[0].checkResult()
  
  // THEN 
  if actual != true {
    t.Error()
  } 
  
  // and GIVEN 
  tests = []testOut{
    {
      "audiofile",
      []result{
        {
          "hh", 
          "good",
        },
        {
          "iy", 
          "good", 
        },
        {
          "v", 
          "good", 
        },
        {
          "er", 
          "surprise", 
        },
      },
      []expectation{
        simpleExpect{
          []string{
            "hh",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "iy",
          },
          []string{
            "good", "possible",
          },        
        },
        simpleExpect{
          []string{
            "v",
          },
          []string{
            "good", "possible",
          },        
        },
      },
      false,
    },    
  }
  
  // WHEN 
  actual = tests[0].checkResult()
  
  // THEN 
  if actual == true {
    t.Error()
  }      
}

func Test_allPass(t *testing.T) {
  // GIVEN 
  orE := orExpect{
    [][]simpleExpect{
      {
        {
          []string{
            "*",
          },
          []string{
            "surprise",
          },
        },
        {
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },
        },
        {
          []string{
            "s",
          },
          []string{
            "good", "possible", "missing",
          },
        },
      },
      {
        {
          []string{
            "s",
          },
          []string{
            "missing",
          },
        },
        {
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },
        },
        {
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },
        },              
      },
      {
        {
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },
        },              
        {
          []string{
            "s",
          },
          []string{
            "missing",
          },
        },
        {
          []string{
            "*",
          },
          []string{
            "surprise", "absent",
          },
        },              
      },
    },
  }
  result := results{
    {
      "ae", "surprise",
    },
    {
      "s", "missing",
    },
    {
      "z", "surprise",
    },
    {
      "k", "good",
    },
    {
      "uw", "good",
    },
    {
      "l", "good",
    },
  }

  // WHEN 
  actual, ok := orE.allPass(result)
  expected := []results{
    {
      {
        "z", "surprise",
      },
      {
        "k", "good",
      },
      {
        "uw", "good",
      },
      {
        "l", "good",
      },
    },
    {
      {
        "k", "good",
      },
      {
        "uw", "good",
      },
      {
        "l", "good",
      },      
    },
  }
  
  // THEN 
  if !ok {
    t.Error()
  } else if !reflect.DeepEqual(actual, expected){
    fmt.Println("actual =", actual)
    t.Error()
  }
}