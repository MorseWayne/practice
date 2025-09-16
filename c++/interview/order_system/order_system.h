#include <iostream>
#include <string>
#include <unordered_map>

using namespace std;

enum OrderStatus { UNDEFINED, PENDING = 1, REJECTED, CONFIRMED, COMPLETED };

struct OrderRequest {
    string order_id;
    string account;
    string item_id;
    string side;  // BUY/SELL
    double price;
    int quantity;
};

struct OrderUpdate {
    string order_id;
    OrderStatus status;
};

class OrderSystem {
public:
    bool Save(const OrderRequest& req)
    {
        // 冲突检测：同一用户+同一商品；未终结订单（PENDING/CONFIRMED）之间买价 >= 卖价
        for (const auto& kv : orders_) {
            const OrderRequest& exist = kv.second;
            auto itStatus = ordersStatus_.find(exist.order_id);
            if (itStatus == ordersStatus_.end())
                continue;
            OrderStatus s = itStatus->second;
            if (s != PENDING && s != CONFIRMED)
                continue;  // 仅未终结订单参与冲突检测

            // 检查同账户 同商品下 是否违背买卖规则： 卖出价格不能低于买入价格，违背该规则，拒绝该订单
            if (exist.account == req.account && exist.item_id == req.item_id) {
                if (req.side == "BUY" && exist.side == "SELL") {
                    if (req.price >= exist.price)
                        return false;
                }
                if (req.side == "SELL" && exist.side == "BUY") {
                    if (exist.price >= req.price)
                        return false;
                }
            }
        }

        // 保存新订单为待处理状态
        orders_[req.order_id] = req;
        ordersStatus_[req.order_id] = CONFIRMED;
        return true;
    }

    bool Update(const OrderUpdate& update)
    {
        if (orders_.count(update.order_id) == 0) {
            return false;
        }

        OrderStatus curr = ordersStatus_[update.order_id];

        // 合法状态流转：
        // PENDING -> REJECTED | CONFIRMED
        // CONFIRMED -> COMPLETED
        if (curr == PENDING && (update.status == REJECTED || update.status == CONFIRMED)) {
            ordersStatus_[update.order_id] = update.status;
            return true;
        }
        if (curr == CONFIRMED && update.status == COMPLETED) {
            ordersStatus_[update.order_id] = update.status;
            return true;
        }
        return false;
    }

private:
    unordered_map<string, OrderRequest> orders_;
    unordered_map<string, OrderStatus> ordersStatus_;
};

OrderSystem g_system;

// 1. 买卖订单同时处于未完成状态：PENDING，CONFIRMED
// 2. 卖出价格低于买入价格
bool checkAndSaveOrder(const OrderRequest& req) { return g_system.Save(req); }

bool updateOrderStatus(const OrderUpdate& update) { return g_system.Update(update); }