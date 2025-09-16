#include <catch2/catch.hpp>
#include <iostream>

using namespace std;
class A {
    virtual void test() {
        cout << "A: tset" << endl;        
    }
};

class B : virtual public A {
    void test() override {
        cout << "B: test" << endl;
    }
};

class C : virtual public A {
    virtual void test() {
        cout << "C: test" << endl;
    }
};

class D : virtual public B, virtual public C {
    virtual void test() {
        cout << "C: test" << endl;
    }
};


TEST_CASE("Virtual table basic test", "[vtable]")
{
    D d;
}