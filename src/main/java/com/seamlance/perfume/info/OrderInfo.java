package com.seamlance.perfume.info;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

import java.math.BigDecimal;
import java.util.List;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class OrderInfo extends BaseInfo {
    private BigDecimal totalPrice;
    private String status;
    private UserInfo userInfo;

    List<OrderItemInfo> orderItemInfoList;

    public BigDecimal getTotalPrice() {
        return totalPrice;
    }

    public void setTotalPrice(BigDecimal totalPrice) {
        this.totalPrice = totalPrice;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public UserInfo getUserInfo() {
        return userInfo;
    }

    public void setUserInfo(UserInfo userInfo) {
        this.userInfo = userInfo;
    }

    public List<OrderItemInfo> getOrderItemInfoList() {
        return orderItemInfoList;
    }

    public void setOrderItemInfoList(List<OrderItemInfo> orderItemInfoList) {
        this.orderItemInfoList = orderItemInfoList;
    }
}
