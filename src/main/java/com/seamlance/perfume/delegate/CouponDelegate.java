package com.seamlance.perfume.delegate;

import java.util.List;

import com.seamlance.perfume.info.CouponInfo;

public interface CouponDelegate {
    CouponInfo createCoupon(CouponInfo couponInfo) throws Exception;
    CouponInfo updateCoupon(String id, CouponInfo couponInfo) throws Exception;
    void deleteCoupon(String id) throws Exception;
    CouponInfo getCouponById(String id) throws Exception;
    List<CouponInfo> getAllCoupon() throws Exception;
}
