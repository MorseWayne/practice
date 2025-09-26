#include <catch2/catch.hpp>
#include <iostream>

using namespace std;
class A {
public:
    ~A() { cout << "~A: called" << endl; }
    virtual void test() { cout << "A: test" << endl; }
};

class B : public A {
public:
    virtual ~B() { cout << "~B: called" << endl; }
    void test() override { cout << "B: test" << endl; }
};

#include <memory>
TEST_CASE("Virtual table basic test", "[vtable]")
{
    // std::shared_ptr<B>  a = std::make_shared<B>();
    // a->test();

    A* a1 = new B();
    delete a1;
}