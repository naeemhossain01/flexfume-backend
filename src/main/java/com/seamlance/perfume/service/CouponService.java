package com.seamlance.perfume.service;

import java.util.List;

import com.seamlance.perfume.entity.Coupon;

public interface CouponService {
    Coupon createCoupon(Coupon coupon) throws Exception;
    Coupon updateCoupon(String couponId, Coupon coupon) throws Exception;
    void deleteCoupon(String couponId) throws Exception;
    Coupon getCouponById(String couponId) throws Exception;
    List<Coupon> getAllCoupon() throws Exception;
    Coupon getCouponByCode(String code);
}
