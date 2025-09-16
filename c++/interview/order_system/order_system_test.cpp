#include <catch2/catch.hpp>
#include "order_system.h"

TEST_CASE("system test", "[order_system]")
{
    OrderRequest order1 = {"1", "U1", "A", "BUY", 100, 10};
    REQUIRE(checkAndSaveOrder(order1));

    OrderRequest order2 = {"2", "U1", "A", "SELL", 102, 5};
    REQUIRE(checkAndSaveOrder(order2));

    OrderRequest order3 = {"3", "U1", "A", "SELL", 99, 5};
    REQUIRE(checkAndSaveOrder(order3) == false);

    OrderUpdate update1 = {"1", COMPLETED};
    REQUIRE(updateOrderStatus(update1));

    OrderUpdate update2 = {"1", REJECTED};
    REQUIRE(updateOrderStatus(update2) == false);

    OrderRequest order2_again = {"4", "U1", "A", "SELL", 98, 5};
    REQUIRE(checkAndSaveOrder(order2_again));
}


