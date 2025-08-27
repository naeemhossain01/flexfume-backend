package com.seamlance.perfume.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.seamlance.perfume.info.CartInfo;
import com.seamlance.perfume.info.CouponInfo;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class CouponRequest {

    private List<String> cartInfoList;
    private String couponCode;

    public String getCouponCode() {
        return couponCode;
    }

    public void setCouponCode(String couponCode) {
        this.couponCode = couponCode;
    }

    public List<String> getCartInfoList() {
        return cartInfoList;
    }

    public void setCartInfoList(List<String> cartInfoList) {
        this.cartInfoList = cartInfoList;
    }
}
