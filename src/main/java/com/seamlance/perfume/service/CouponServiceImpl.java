package com.seamlance.perfume.service;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

import com.seamlance.perfume.dao.CouponDao;
import com.seamlance.perfume.dao.CouponUsageDao;
import com.seamlance.perfume.entity.Coupon;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.exception.ResourceNotFoundException;


@Service
public class CouponServiceImpl implements CouponService {

    @Autowired
    private CouponDao couponDao;
    
    @Autowired
    private CouponUsageDao couponUsageDao;

    @Override
    public Coupon createCoupon(Coupon coupon) throws Exception {
        if(coupon == null) {
            throw new InvalidRequestsException("Invalid Coupon");
        }

        try {
            coupon = couponDao.saveAndFlush(coupon);
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return coupon;
    }

    @Override
    public Coupon updateCoupon(String couponId, Coupon coupon) throws Exception {
        Coupon updatedCoupon = couponDao.findById(couponId).orElseThrow(() -> new ResourceNotFoundException("Coupon not found"));

        try {
            if(coupon.getCouponType() != null) updatedCoupon.setCouponType(coupon.getCouponType());
            if(coupon.getAmount() != null) updatedCoupon.setAmount(coupon.getAmount());
            if(coupon.getUsageLimit() != 0 ) updatedCoupon.setUsageLimit(coupon.getUsageLimit());
            if(coupon.getMinOrderAmount() != null) updatedCoupon.setMinOrderAmount(coupon.getMinOrderAmount());
            if(coupon.getExpirationTime() != null) updatedCoupon.setExpirationTime(coupon.getExpirationTime());
            if(coupon.getMaxAmountApplied() != null) updatedCoupon.setMaxAmountApplied(coupon.getMaxAmountApplied());
            if(coupon.isActive() != updatedCoupon.isActive()) updatedCoupon.setActive(coupon.isActive());
            if(coupon.getCode() != null) updatedCoupon.setCode(coupon.getCode());

            updatedCoupon = couponDao.save(updatedCoupon);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return updatedCoupon;
    }

    @Override
    public void deleteCoupon(String couponId) throws Exception {
        Coupon coupon = couponDao.findById(couponId).orElseThrow(() -> new ResourceNotFoundException("Coupon not found"));
        
        try {
            // Check if there are any usage records for this coupon
            Integer usageCount = couponUsageDao.getTotalUsageCountByCouponId(couponId);
            if (usageCount != null && usageCount > 0) {
                throw new InvalidRequestsException("Cannot delete coupon that has been used. Found " + usageCount + " usage records.");
            }
            
            couponDao.delete(coupon);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException("Failed to delete coupon: " + e.getMessage());
        }
    }

    @Override
    public Coupon getCouponById(String couponId) {
        return couponDao.findById(couponId).orElseThrow(() -> new ResourceNotFoundException("Coupon not found"));
    }

    @Override
    public List<Coupon> getAllCoupon() {
        List<Coupon> coupons = couponDao.findAll();
        
        // Get usage statistics for all coupons efficiently
        if (!coupons.isEmpty()) {
            List<String> couponIds = coupons.stream()
                .map(Coupon::getId)
                .collect(Collectors.toList());
            
            // Get batch usage statistics
            Object[][] usageStats = couponUsageDao.getUsageStatisticsByCouponIds(couponIds);
            
            // Create a map for quick lookup
            Map<String, Object[]> statsMap = new HashMap<>();
            for (Object[] stat : usageStats) {
                String couponId = (String) stat[0];
                statsMap.put(couponId, stat);
            }
            
            // Note: Since Coupon entity doesn't have usage fields, 
            // this aggregation would be better handled in the mapper/delegate layer
            // when converting to CouponInfo
        }
        
        return coupons;
    }

    @Override
    public Coupon getCouponByCode(String code) {

        if(code == null) {
            throw new ResourceNotFoundException("Coupon not found");
        }

        Coupon coupon =  null;

        try {
            coupon = couponDao.findByCode(code);
        } catch (Exception e) {
            throw new ResourceNotFoundException(e.getMessage());
        }

        return coupon;
    }
}
