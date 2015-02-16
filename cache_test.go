package cache

import (
	"testing"
	"time"
    "reflect"
)

func TestOK_SizeOf_String(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	obj := "0123456789";
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 10;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_Int64(t *testing.T) {
	// enable logger
	EnableLogger(true)
	
	// test data
	var obj int64 = 100;
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 8;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_Int32(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	var obj int32 = 100;
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 4;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_MapByte(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	obj := make(map[int]byte, 10);
	for i := 0; i < 10; i++ {
		obj[i] = byte(1);
	}
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 10;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_MapInt32(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	obj := make(map[int]int32, 10);
	for i := 0; i < 10; i++ {
		obj[i] = 1;
	}
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 40;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_SliceByte(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	obj := make([]byte, 10);
	for i := range obj {
		obj[i] = byte(1);
	}
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 10;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_SliceInt32(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	obj := make([]int32, 10);
	for i := range obj {
		obj[i] = 1;
	}
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 40;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_ArrayByte(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	var obj = [10]byte{};
	for i := range obj {
		obj[i] = byte(1);
	}
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 10;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_SizeOf_ArrayInt32(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	obj := [10]int32{};
	for i := range obj {
		obj[i] = 1;
	}
	
	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// SizeOf
	expected := 40;
	size := c.SizeOf(obj)
	if size != expected {
		t.Errorf("size(%d) is invalid. expected = %d", size, expected)
	}
}

func TestOK_Priority_AccessCount(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	key1 := "key1"
	key2 := "key2"
	key3 := "key3"
	type Obj struct {
		v1 string
		v2 int
	}
	value1 := &Obj{v1: "test1", v2: 1,}
	value2 := &Obj{v1: "test2", v2: 2,}
	value3 := &Obj{v1: "test3", v2: 3,}

	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 2, // Access count of 2 or more, priority +1
		}
	c := New(opt)
	
	// set
	c.Set(key1, value1, time.Duration(-1)) //expiration
	c.Set(key2, value2, time.Duration(-1)) //expiration
	c.Set(key3, value3, time.Duration(-1)) //expiration
	
	//check data
	i1, _ := c.GetItem(key1)
	if i1.Priority != 1 {
		t.Errorf("priority(%d) is invalid. expected = %d", i1.Priority, 1)
	}
	i2, _ := c.GetItem(key2)
	if i2.Priority != 1 {
		t.Errorf("priority(%d) is invalid. expected = %d", i2.Priority, 1)
	}
	
	// Priority up
	_, _ = c.Get(key1) //access count up 1
	_, _ = c.Get(key1) //access count up 2
	_, _ = c.Get(key2) //access count up 1
	c.Optimize()
	
	//check data
	i1, found := c.GetItem(key1)
	if found {
		if i1.Priority != 2 {
			t.Errorf("priority(%d) is invalid. expected = %d", i1.Priority, 2)
		}
		i2, _ = c.GetItem(key2)
		if i2.Priority != 1 {
			t.Errorf("priority(%d) is invalid. expected = %d", i2.Priority, 1)
		}
	} else {
		t.Errorf("item is not found.")
	}
}

func TestOK_Priority_Expiration(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	key1 := "key1"
	key2 := "key2"
	type Obj struct {
		v1 string
		v2 int
	}
	value1 := &Obj{v1: "test1", v2: 1,}
	value2 := &Obj{v1: "test2", v2: 2,}

	// init cache
	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	
	// set
	c.Set(key1, value1, time.Duration(-1)) //expiration
	c.Set(key2, value2, time.Duration(1000)) //expiration
	
	//check data
	i1, _ := c.GetItem(key1)
	if i1.Priority != 1 {
		t.Errorf("priority(%v) is invalid. expected = %v", i1.Priority, 1)
	}
	i2, _ := c.GetItem(key2)
	if i2.Priority != 1 {
		t.Errorf("priority(%v) is invalid. expected = %v", i2.Priority, 1)
	}
	time.Sleep(1000)
	if c.Priority(key2) != 0 { //priority will be 0 by expiration 
		t.Errorf("priority(%v) is invalid. expected = %v", i2.Priority, 1)
	}
	
	// Priority up
	c.Optimize() // i2 is deleted by expiration.
	
	//check data
	i1, found := c.GetItem(key1)
	if found {
		if i1.Priority != 1 {
			t.Errorf("priority(%v) is invalid. expected = %v", i1.Priority, 1)
		}
		_, found := c.GetItem(key2)
		if found {
			t.Errorf("item is found. expected = %v", false) 
		}
	} else {
		t.Errorf("item is not found.")
	}
}

func TestOK_Priority_Access(t *testing.T) {
	// enable logger
	EnableLogger(true)

	// test data
	key1 := "key1"
	key2 := "key2"
	type Obj struct {
		V1 string
		V2 int
	}
	value1 := &Obj{V1: "test1", V2: 1,}
	value2 := &Obj{V1: "test2", V2: 2,}

	// init cache
	opt := Option{
			ThresholdAccess: 10000000000, //Access of 10 seconds or less, priority +1
			ThresholdAccessCount: 10, // Access count of 10 or more, priority +1
		}
	c := New(opt)
	
	// set
	c.Set(key1, value1, time.Duration(-1)) //expiration
	c.Set(key2, value2, time.Duration(-1)) //expiration
	
	//check data
	i1, _ := c.GetItem(key1)
	if i1.Priority != 1 {
		t.Errorf("priority(%v) is invalid. expected = %v", i1.Priority, 1)
	}
	i2, _ := c.GetItem(key2)
	if i2.Priority != 1 {
		t.Errorf("priority(%v) is invalid. expected = %v", i2.Priority, 1)
	}
	
	// Priority up
	_, _ = c.Get(key1) //update last access time
	//_, _ = c.Get(key2) //update last access time
	c.Optimize()
	
	//check data
	i1, found := c.GetItem(key1)
	if found {
		if i1.Priority != 2 {
			t.Errorf("priority(%v) is invalid. expected = %v", i1.Priority, 2)
		}
		i2, _ = c.GetItem(key2)
		if i2.Priority != 1 {
			t.Errorf("priority(%v) is invalid. expected = %v", i2.Priority, 1)
		}
	} else {
		t.Errorf("item is not found.")
	}
}

func TestOK_SetGet_OverWrite(t *testing.T) {
	// enable logger
	EnableLogger(true)

	key := "TestOK_Set_SameObject"
	type Obj struct {
		v1 string
		v2 int
	}
	value1 := &Obj{v1: "test1", v2: 1,}
	value2 := &Obj{v1: "test1", v2: 1,}

	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	
	c.Set(key, value1, time.Duration(10))
	c.Set(key, value2, time.Duration(10))
	v, found := c.Get(key)
	if !found {
		t.Errorf("key(%v) is not found.", key)
	}
	if v == nil {
		t.Errorf("value of key(%v) is nil.", key)
		t.FailNow()
	}
	if reflect.ValueOf(*v).Type() != reflect.ValueOf(value1).Type() {
		t.Errorf("v(%v) is not same type with value(%v).", reflect.ValueOf(*v).Type(), reflect.ValueOf(value1).Type())
		t.FailNow()
	}
	if *v != value1 {
		t.Logf("v(%v) is not same value with value1(%v).", v, value1)
	}
	if *v != value2 {
		t.Errorf("v(%v) is not same value with value2(%v).", v, value2)
		t.FailNow()
	}
}

func TestOK_SetGet_String(t *testing.T) {
	// enable logger
	EnableLogger(true)

	key := "TestSet_string"
	value := "testString"

	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	
	c.Set(key, value, time.Duration(10))
	v, found := c.Get(key)
	if !found {
		t.Errorf("key(%v) is not found.", key)
	}
	if v == nil {
		t.Errorf("value of key(%v) is nil.", key)
		t.FailNow()
	}
	if reflect.ValueOf(*v).Type() != reflect.ValueOf(value).Type() {
		t.Errorf("v(%v) is not same type with value(%v).", reflect.ValueOf(*v).Type(), reflect.ValueOf(value).Type())
		t.FailNow()
	}
	if *v != value {
		t.Errorf("v(%v) is not same value with value(%v).", v, value)
		t.FailNow()
	}
}

func TestNG_SetGet_String(t *testing.T) {
	// enable logger
	EnableLogger(true)

	key := "TestSet_string"
	badKey := "badKey"
	value := "testString"

	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	c.Set(key, value, time.Duration(10))
	v, found := c.Get(badKey)
	if !found {
		t.Logf("key(%v) is not found.", key)
		return
	}
	if v == nil {
		t.Errorf("value of key(%v) is nil.", key)
		t.FailNow()
	}
	if reflect.ValueOf(v).Type() != reflect.ValueOf(value).Type() {
		t.Errorf("v(%v) is not same type with value(%v).", reflect.ValueOf(v).Type(), reflect.ValueOf(value).Type())
		t.FailNow()
	}
	if *v != value {
		t.Errorf("v(%v) is not same value with value(%v).", v, value)
		t.FailNow()
	}
}

func TestNG_Set_Chan(t *testing.T) {
	// enable logger
	EnableLogger(true)

	key := "TestSet_chan"
	value := make(chan int, 1)

	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	err := c.Set(key, value, time.Duration(10))
	if (err != nil) {
		t.Logf("expected error. error = %s", err.Error())
		return
	}
	t.FailNow()
}

func TestNG_Set_Func(t *testing.T) {
	// enable logger
	EnableLogger(true)

	key := "TestSet_Func"
	value := func(i int) int { return i+1 }

	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	err := c.Set(key, value, time.Duration(10))
	if (err != nil) {
		t.Logf("expected error. error = %s", err.Error())
		return
	}
	t.FailNow()
}

func TestOK_Del_String(t *testing.T) {
	// enable logger
	EnableLogger(true)

	key := "TestOK_Del_String"
	value := "testString"

	opt := Option{
			ThresholdAccess: 0,
			ThresholdAccessCount: 0,
		}
	c := New(opt)
	
	c.Set(key, value, time.Duration(10))
	v, found := c.Get(key)
	if !found {
		t.Errorf("key(%v) is not found.", key)
		t.FailNow()
	}
	if v == nil {
		t.Errorf("value of key(%v) is nil.", key)
		t.FailNow()
	}
	if reflect.ValueOf(*v).Type() != reflect.ValueOf(value).Type() {
		t.Errorf("v(%v) is not same type with value(%v).", reflect.ValueOf(*v).Type(), reflect.ValueOf(value).Type())
		t.FailNow()
	}
	if *v != value {
		t.Errorf("v(%v) is not same value with value(%v).", v, value)
		t.FailNow()
	}
	
	//delete
	c.Del(key)
	v, found = c.Get(key)
	if found {
		t.Errorf("key(%v) is found.", key)
		t.FailNow()
	}
}
