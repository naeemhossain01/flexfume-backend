package com.seamlance.perfume.controller;


import java.math.BigDecimal;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;

import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.delegate.CouponUsageDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.dto.CouponRequest;
import com.seamlance.perfume.exception.InvalidRequestsException;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/coupon-usage")
public class CouponUsageController {
    @Autowired
    private CouponUsageDelegate couponUsageDelegate;

    @PostMapping()
    @PreAuthorize("hasAuthority('USER') or hasAuthority('ADMIN')")
    public ApiResponse<BigDecimal> applyCoupon(@RequestBody CouponRequest couponRequest) throws InvalidRequestsException {
        return ApiResponse.success(couponUsageDelegate.applyCoupon(couponRequest.getCartInfoList(), couponRequest.getCouponCode()));
    }

    @DeleteMapping()
    public ApiResponse<String> deleteCouponUsage(@RequestParam String code) {
        couponUsageDelegate.deleteCouponUsage(code);
        return ApiResponse.success(Constant.COUPON_REMOVED);
    }
}
