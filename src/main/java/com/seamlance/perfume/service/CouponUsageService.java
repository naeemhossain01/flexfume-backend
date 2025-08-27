package com.seamlance.perfume.service;

import com.seamlance.perfume.entity.Cart;
import com.seamlance.perfume.entity.CouponUsage;
import com.seamlance.perfume.exception.InvalidRequestsException;

import java.math.BigDecimal;
import java.util.List;

public interface CouponUsageService {
    BigDecimal applyCoupon(List<String> cartIdList, String couponCode) throws InvalidRequestsException;
    CouponUsage getCouponUsage(String code, String userId);
    void deleteCouponUsage(String couponCode);
}
