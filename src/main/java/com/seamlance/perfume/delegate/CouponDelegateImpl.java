package com.seamlance.perfume.delegate;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import com.seamlance.perfume.dao.CouponUsageDao;
import com.seamlance.perfume.entity.Coupon;
import com.seamlance.perfume.info.CouponInfo;
import com.seamlance.perfume.mapper.CouponMapper;
import com.seamlance.perfume.service.CouponService;

@Component
public class CouponDelegateImpl implements CouponDelegate {

    @Autowired
    private CouponService couponService;

    @Autowired
    private CouponMapper couponMapper;
    
    @Autowired
    private CouponUsageDao couponUsageDao;

    @Override
    public CouponInfo createCoupon(CouponInfo couponInfo) throws Exception {
        return couponMapper.toInfo(couponService.createCoupon(couponMapper.toEntity(couponInfo)));
    }

    @Override
    public CouponInfo updateCoupon(String id, CouponInfo couponInfo) throws Exception {
        return couponMapper.toInfo(couponService.updateCoupon(id, couponMapper.toEntity(couponInfo)));
    }

    @Override
    public void deleteCoupon(String id) throws Exception {
        couponService.deleteCoupon(id);
    }

    @Override
    public CouponInfo getCouponById(String id) throws Exception {
        CouponInfo couponInfo = couponMapper.toInfo(couponService.getCouponById(id));
        
        // Enrich with usage statistics for single coupon
        enrichWithUsageStatistics(couponInfo, id);
        
        return couponInfo;
    }

    @Override
    public List<CouponInfo> getAllCoupon() throws Exception {
        List<Coupon> coupons = couponService.getAllCoupon();
        List<CouponInfo> couponInfos = couponMapper.toInfoList(coupons);
        
        // Enrich with usage statistics for all coupons efficiently
        enrichWithUsageStatistics(couponInfos);
        
        return couponInfos;
    }
    
    /**
     * Enrich a single CouponInfo with usage statistics
     */
    private void enrichWithUsageStatistics(CouponInfo couponInfo, String couponId) {
        Integer totalUsage = couponUsageDao.getTotalUsageCountByCouponId(couponId);
        BigDecimal totalSavings = couponUsageDao.getTotalSavingsAmountByCouponId(couponId);
        Integer uniqueUsers = couponUsageDao.getUniqueUsersCountByCouponId(couponId);
        
        couponInfo.setTotalUsageCount(totalUsage != null ? totalUsage : 0);
        couponInfo.setTotalSavingsAmount(totalSavings != null ? totalSavings : BigDecimal.ZERO);
        couponInfo.setUniqueUsersCount(uniqueUsers != null ? uniqueUsers : 0);
    }
    
    /**
     * Enrich multiple CouponInfos with usage statistics efficiently using batch query
     */
    private void enrichWithUsageStatistics(List<CouponInfo> couponInfos) {
        if (couponInfos.isEmpty()) {
            return;
        }
        
        // Get all coupon IDs
        List<String> couponIds = couponInfos.stream()
            .map(CouponInfo::getId)
            .collect(Collectors.toList());
        
        // Get batch usage statistics
        Object[][] usageStats = couponUsageDao.getUsageStatisticsByCouponIds(couponIds);
        
        // Create a map for quick lookup: couponId -> [totalUsage, totalSavings, uniqueUsers]
        Map<String, Object[]> statsMap = new HashMap<>();
        for (Object[] stat : usageStats) {
            String couponId = (String) stat[0];
            statsMap.put(couponId, stat);
        }
        
        // Enrich each CouponInfo
        for (CouponInfo couponInfo : couponInfos) {
            Object[] stats = statsMap.get(couponInfo.getId());
            if (stats != null) {
                // stats[0] = couponId, stats[1] = totalUsage, stats[2] = totalSavings, stats[3] = uniqueUsers
                Long totalUsage = (Long) stats[1];
                BigDecimal totalSavings = (BigDecimal) stats[2];
                Long uniqueUsers = (Long) stats[3];
                
                couponInfo.setTotalUsageCount(totalUsage != null ? totalUsage.intValue() : 0);
                couponInfo.setTotalSavingsAmount(totalSavings != null ? totalSavings : BigDecimal.ZERO);
                couponInfo.setUniqueUsersCount(uniqueUsers != null ? uniqueUsers.intValue() : 0);
            } else {
                // No usage data for this coupon
                couponInfo.setTotalUsageCount(0);
                couponInfo.setTotalSavingsAmount(BigDecimal.ZERO);
                couponInfo.setUniqueUsersCount(0);
            }
        }
    }
}
