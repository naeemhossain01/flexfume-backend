package com.seamlance.perfume.controller;


import com.seamlance.perfume.delegate.CouponDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.CouponInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/coupon")
public class CouponController {
    @Autowired
    private CouponDelegate couponDelegate;

    @PostMapping()
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<CouponInfo> create(@RequestBody CouponInfo couponInfo) throws Exception {
        return ApiResponse.success(couponDelegate.createCoupon(couponInfo));
    }

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<CouponInfo> update(@PathVariable String id, @RequestBody CouponInfo couponInfo) throws Exception {
        return ApiResponse.success(couponDelegate.updateCoupon(id, couponInfo));
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<String> delete(@PathVariable String id) throws Exception {
        couponDelegate.deleteCoupon(id);
        return ApiResponse.success("Coupon deleted successfully");
    }

    @GetMapping("/{id}")
    public ApiResponse<CouponInfo> get(@PathVariable String id) throws Exception {
        return ApiResponse.success(couponDelegate.getCouponById(id));
    }

    @GetMapping("/all")
    public ApiResponse<List<CouponInfo>> getAll() throws Exception {
        return ApiResponse.success(couponDelegate.getAllCoupon());
    }
}
