// MIT License
//
// Copyright (c) 2016 GenesisCommunity
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package script

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

type TestVM struct {
	Input  string
	Func   string
	Output string
}

func (block *Block) String() (ret string) {
	if (*block).Objects != nil {
		ret = fmt.Sprintf("Objects: %v\n", (*block).Objects)
	}
	ret += fmt.Sprintf("Type: %v \n", (*block).Type)
	if (*block).Children != nil {
		ret += fmt.Sprintf("Blocks: [\n")
		for i, item := range (*block).Children {
			ret += fmt.Sprintf("{%d: %v}\n", i, item.String())
		}
		ret += fmt.Sprintf("]")
	}
	return
}

func getMap() map[string]interface{} {
	return map[string]interface{}{`par0`: `Parameter 0`, `par1`: `Parameter 1`}
}

func getArray() []interface{} {
	return []interface{}{map[string]interface{}{`par0`: `Parameter 0`, `par1`: `Parameter 1`},
		"The second string", int64(2000)}
}

// Str converts the value to a string
func str(v interface{}) (ret string) {
	return fmt.Sprint(v)
}

func lenArray(par []interface{}) int64 {
	return int64(len(par))
}

func Money(v interface{}) (ret decimal.Decimal) {
	ret, _ = ValueToDecimal(v)
	return ret
}

func TestVMCompile(t *testing.T) {
	test := []TestVM{
		{`contract sets {
			settings {
				val = 1.56
				rate = 100000000000
				name="Name parameter"
			}
			action {
				$result = Settings("@22sets","name")
			}
		}
		func result() string {
			var par map
			return CallContract("@22sets", par) + "=" + sets()
		}
		`, `result`, `Name parameter=Name parameter`},

		{`func proc(par string) string {
					return par + "proc"
					}
				func forarray string {
					var my map
					var ret array
					var myret array
					
					ret = GetArray()
					myret[1] = "Another "
					my = ret[0]
					my["par3"] = 3456
					ret[2] = "Test"
					return Sprintf("result=%s+%s+%d+%s", ret[1], my["par0"], my["par3"], myret[1] + ret[2])
				}`, `forarray`, `result=The second string+Parameter 0+3456+Another Test`},
		{`func proc(par string) string {
				return par + "proc"
				}
			func formap string {
				var my map
				var ret map

				ret = GetMap()
		//			Println(ret)
				//Println("Ooops", ret["par0"], ret["par1"])
				my["par1"] = "my value" + proc(" space ")
				my["par2"] = 203 * (100-86)
				return Sprintf("result=%s+%d+%s+%s+%d", ret["par1"], my["par2"] + 32, my["par1"], proc($glob["test"] ), $glob["number"] )
			}`, `formap`, `result=Parameter 1+2874+my value space proc+String valueproc+1001`},
		{`func runtime string {
						var i int
						i = 50
						return Sprintf("val=%d", i 0)
					}`, `runtime`, `runtime panic error`},
		{`func nop {
							return
						}

						func loop string {
							var i int
							while true {//i < 10 {
								i=i+1
								if i==5 {
									continue
								}
								if i == 121 {
									i = i+ 4
									break
								}
							}
							nop()
							return Sprintf("val=%d", i)
						}`, `loop`, `val=125`},
		{`contract my {
							data {
								Par1 int
								Par2 string
							}
							func conditions {
								var q int
								Println("Front", $Par1, $parent)
				//				my("Par1,Par2,ext", 123, "Parameter 2", "extended" )
							}
							func action {
								Println("Main", $Par2, $ext)
							}
						}
						contract mytest {
							func init string {
								empty()
								my("Par1,Par2,ext", 123, "Parameter 2", "extended" )
								//my("Par1,Par2,ext", 33123, "Parameter 332", "33extended" )
								//@26empty("test",10)
								empty("toempty", 10)
								Println( "mytest", $parent)
								return "OK"
							}
						}
						contract empty {
							conditions {Println("EmptyCond")
								}
							action {
								Println("Empty", $parent)
								if 1 {
									my("Par1,Par2,ext", 123, "Parameter 2", "extended" )
								}
							}
						}
						`, `mytest.init`, `OK`},
		{`func line_test string {
						return "Start " +
						Sprintf( "My String %s %d %d",
								"Param 1", 24,
							345 + 789)
					}`, `line_test`, `Start My String Param 1 24 1134`},

		{`func err_test string {
						if 1001.02 {
							error "Error message err_test"
						}
						return "OK"
					}`, `err_test`, `{"type":"error","error":"Error message err_test"}`},
		{`contract my {
					data {
						PublicKey  bytes
						FirstName  string
						MiddleName string "optional"
						LastName   string
					}
					func init string {
						return "OK"
					}
				}`, `my.init`, `OK`},

		{`func temp3 string {
						var i1 i2 int, s1 string, s2 string
						i2, i1 = 348, 7
						if i1 > 5 {
							var i5 int, s3 string
							i5 = 26788
							s1 = "s1 string"
							i2 = (i1+2)*i5+i2
							s2 = Sprintf("temp 3 function %s %d", Sprintf("%s + %d", s1, i2), -1 )
						}
						return s2
					}`, `temp3`, `temp 3 function s1 string + 241440 -1`},
		{`func params2(myval int, mystr string ) string {
					if 101>myval {
						if myval == 90 {
						} else {
							return Sprintf("myval=%d + %s", myval, mystr )
						}
					}
					return "OOPs"
				}
				func temp2 string {
					if true {
						return params2(51, "Params 2 test")
					}
				}
				`, `temp2`, `myval=51 + Params 2 test`},

		{`func params(myval int, mystr string ) string {
					return Sprintf("Params function %d %s", 33 + myval + $test1, mystr + " end" )
				}
				func temp string {
					return "Prefix " + params(20, "Test string " + $test2) + $test3( 202 )
				}
				`, `temp`, `Prefix Params function 154 Test string test 2 endtest=202=test`},
		{`func my_test string {
				return Sprintf("Called my_test %s %d", "Ooops", 777)
			}

	contract my {
			func initf string {
				return Sprintf("%d %s %s %s", 65123 + (1001-500)*11, my_test(), "Test message", Sprintf("> %s %d <","OK", 999 ))
			}
	}`, `my.initf`, `70634 Called my_test Ooops 777 Test message > OK 999 <`},
		{`contract vars {
		func cond() string {return "vars"}
		func actions() { var test int}
	}`, `vars.cond`, `vars`},
		{`func mytail(name string, tail ...) string {
		if lenArray(tail) == 0 {
			return name
		}
		if lenArray(tail) == 1 {
			return Sprintf("%s=%v ", name, tail[0])
		}
		return Sprintf("%s=%v+%v ", name, tail[1], tail[0])
	}
	func emptytail(tail ...) string {
		return Sprintf("%d ", lenArray(tail))
	}
	func sum(out string, values ...) string {
		var i, res int
		while i < lenArray(values) {
		   res = res + values[i]
		   i = i+1
		}
		return Sprintf(out, res)
	}
	func calltail() string {
		var out string
		out = emptytail() + emptytail(10) + emptytail("name1", "name2")
		out = out + mytail("OK") + mytail("1=", 11) + mytail("2=", "name", 11)
		return out + sum("Sum: %d", 10, 20, 30, 40)
	}
	`, `calltail`, `0 1 2 OK1==11 2==11+name Sum: 100`},
		{`func DBFind( table string).Columns(columns string) 
		. Where(format string, tail ...). Limit(limit int).
		Offset(offset int) string  {
		Println("DBFind", table, tail)
		return Sprintf("%s %s %s %d %d=", table, columns, format, limit, offset)
	}
	func names() string {
		var out, cols string
		cols = "name,value"
		out = DBFind( "mytable") + DBFind( "keys"
			).Columns(cols)+ DBFind( "keys"
				).Offset(199).Columns("qq"+"my")
		out = out + DBFind( "table").Columns("name").Where("id=?", 
			100).Limit(10) + DBFind( "table").Where("request")
		return out
	}`, `names`, `mytable   0 0=keys name,value  0 0=keys qqmy  0 199=table name id=? 10 0=table  request 0 0=`},
		{`contract seterr {
				func getset string {
					var i int
					i = MyFunc("qqq", 10)
					return "OK"
				}
			}`, `seterr.getset`, `unknown identifier MyFunc`},
		{`func one() int {
				return 9
			}
			func signfunc string {
				var myarr array
				myarr[0] = 0
				myarr[1] = 1
				var i, k, j int
				k = one()-2
				j = /*comment*/-3
				i = lenArray(myarr) - 1
				return Sprintf("%s %d %d %d %d %d", "ok", lenArray(myarr)-1, i, k, j, -4)
			}`, `signfunc`, `ok 1 1 7 -3 -4`},
		{`func exttest() string {
				return Replace("text", "t")
			}
			`, `exttest`, `function Replace must have 4 parameters`},
		{`func mytest(first string, second int) string {
				return Sprintf("%s %d", first, second)
		}
		func test() {
			return mytest("one", "two")
		}
		`, `test`, `parameter 2 has wrong type`},
		{`func mytest(first string, second int) string {
								return Sprintf("%s %d", first, second)
						}
						func test() string {
							return mytest("one")
						}
						`, `test`, `wrong count of parameters`},
		{
			`func ifMap string {
				var m map
				if m {
					return "empty"
				}
				
				m["test"]=1
				if m {
					return "not empty"
				}

				return error "error"
			}`, "ifMap", "not empty",
		},
		{`func One(list array, name string) string {
			if list {
				var row map 
				row = list[0]
				return row[name]
			}
			return nil
		}
		func Row(list array) map {
			var ret map
			if list {
				ret = list[0]
			}
			return ret
		}
		func GetData().WhereId(id int) array {
			var par array
			var item map
			item["id"] = str(id)
			item["name"] = "Test value " + str(id)
			par[0] = item
			return par
		}
		func GetEmpty().WhereId(id int) array {
			var par array
			return par
		}
		func result() string {
			var m map
			var s string
			m = GetData().WhereId(123).Row()
			s = GetEmpty().WhereId(1).One("name") 
			if s != nil {
				return "problem"
			}
			return m["id"] + "=" + GetData().WhereId(100).One("name")
		}`, `result`, `123=Test value 100`},
		{`func mapbug() string {
			$data[10] = "extend ok"
			return $data[10]
			}`, `mapbug`, `extend ok`},
		{`func result() string {
				var myarr array
				myarr[0] = "string"
				myarr[1] = 7
				myarr[2] = "9th item"
				return Sprintf("RESULT=%s %d %v", myarr...)
			}`, `result`, `RESULT=string 7 9th item`},
		{`func find().Where(pattern string, params ...) string {
				return Sprintf(pattern, params ...)
			}
			func row().Where(pattern string, params ...) string {
				return find().Where(pattern, params ...)
			}
			func result() string {
				return row().Where("%d %d", 10, 20)
			}
			`, `result`, `10 20`},
		{`func result string {
				var arr array
				var mymap map
				arr[100000] = 0
				var i int
				while i < 100 {
					mymap[str(i)] = 10
					i = i + 1
				}
				i = i + "2" 
				i = (i - "10")/"2"*"3"
				return Sprintf("%T %[1]v", .21 + i)
			  }`, `result`, `float64 138.21`},
		{`func money_test string {
				var my2, m1 money
				my2 = 100
				m1 = 1.2
				return Sprintf( "Account %v %v %v", my2/Money(3),  my2 - Money(5.6), m1*Money(5) + Money(my2))
			}`, `money_test`, `Account 33 95 105`},
		{`func long() int {
				return  99999999999999999999
				}
				func result() string {
					return Sprintf("ok=%d", long())
					}`, `result`, `strconv.ParseInt: parsing "99999999999999999999": value out of range 99999999999999999999 [Ln:2 Col:34]`},
		{`func result() string {
			var i, result int
			
			if true {
				if false {
					result = 99
				} else {
					result = 5
				}
			}
			if i == 1 {
				result = 20
			} elif i> 0 {
				result = 30
			} 
			elif i == 0 
			{
				result = result + 50
				if true {
					i=10
				}
			} elif i==10 {
				Println("3")
				result = 0
				i=33
			} elif false {
				Println("4")
				result = 1
			} 
			else 
			{
				Println("5")
				result = 2
			}
			if i == 4 {
				result = result
			} elif i == 20 {
				result = 22
			} else {
				result = result + 23
				i = 11
			}
			if i == 11 {
				result = result + 7
			} else {
				result = 0
			}
			if result == 85 {
				if false {
					result = 1
				} elif 0 {
					result = 5
				} elif 1 {
					result = result + 10
				}
			}
			if result == 10 {
				result = 11
			} elif result == 95 {
				result = result + 1
				if false {
					result = 0
				} elif true {
					result = result + 4
				}
			}
			return Sprintf("%d", result)
		}
					`, `result`, `100`},
	}
	vm := NewVM()
	vm.Extern = true
	vm.Extend(&ExtendData{map[string]interface{}{"Println": fmt.Println, "Sprintf": fmt.Sprintf,
		"GetMap": getMap, "GetArray": getArray, "lenArray": lenArray,
		"str": str, "Money": Money, "Replace": strings.Replace}, nil})

	for ikey, item := range test {
		source := []rune(item.Input)
		if err := vm.Compile(source, &OwnerInfo{StateID: uint32(ikey) + 22, Active: true, TableID: 1}); err != nil {
			if err.Error() != item.Output {
				t.Error(err)
				break
			}
		} else {
			if out, err := vm.Call(item.Func, nil, &map[string]interface{}{
				`rt_state`: uint32(ikey) + 22, `data`: make([]interface{}, 0),
				`test1`: 101, `test2`: `test 2`,
				"glob": map[string]interface{}{`test`: `String value`, `number`: 1001},
				`test3`: func(param int64) string {
					return fmt.Sprintf("test=%d=test", param)
				},
			}); err == nil {
				if out[0].(string) != item.Output {
					t.Error(`error vm ` + out[0].(string) + `!=` + item.Output)
					break
				}
			} else if err.Error() != item.Output {
				t.Error(err)
				break
			}

		}
	}
}

func TestContractList(t *testing.T) {
	test := []TestLexem{{`contract NewContract {
		conditions {
			ValidateCondition($Conditions,$ecosystem_id)
			while i < Len(list) {
				if IsObject(list[i], $ecosystem_id) {
					warning Sprintf("Contract or function %s exists", list[i] )
				}
			}
		}
		action {
		}
		func price() int {
			return  SysParamInt("contract_price")
		}
	}func MyFunc {}`,
		`NewContract,MyFunc`},
		{`contract demo_??ontract {
			data {
				contract_txt str
			}
			func test() {
			}
			conditions {
				if $contract_txt="" {
					warning "Sorry, you do not have contract access to this action."
				}
			}
		} contract another_contract {} func main { func subfunc(){}}`,
			`demo_??ontract,another_contract,main`},
	}
	for _, item := range test {
		list, _ := ContractsList(item.Input)
		if strings.Join(list, `,`) != item.Output {
			t.Error(`wrong names`, strings.Join(list, `,`))
			break
		}
	}
}
