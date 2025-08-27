package com.seamlance.perfume.delegate;

import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.info.CartInfo;

import java.math.BigDecimal;
import java.util.List;

public interface CouponUsageDelegate {
    BigDecimal applyCoupon(List<String> cartInfoList, String couponCode) throws InvalidRequestsException;
    void deleteCouponUsage(String code);
}
