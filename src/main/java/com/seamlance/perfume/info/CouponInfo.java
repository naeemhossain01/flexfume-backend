package com.seamlance.perfume.info;

import java.math.BigDecimal;
import java.time.LocalDateTime;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class CouponInfo extends BaseInfo {
    private String couponType;
    private String code;
    private LocalDateTime expirationTime;
    private int usageLimit;
    private BigDecimal amount;
    private BigDecimal minOrderAmount;
    private BigDecimal maxAmountApplied;
    private boolean active;
    
    // Aggregated usage statistics (calculated, not stored)
    private Integer totalUsageCount;
    private BigDecimal totalSavingsAmount;
    private Integer uniqueUsersCount;

    public String getCouponType() {
        return couponType;
    }

    public void setCouponType(String couponType) {
        this.couponType = couponType;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public LocalDateTime getExpirationTime() {
        return expirationTime;
    }

    public void setExpirationTime(LocalDateTime expirationTime) {
        this.expirationTime = expirationTime;
    }

    public int getUsageLimit() {
        return usageLimit;
    }

    public void setUsageLimit(int usageLimit) {
        this.usageLimit = usageLimit;
    }

    public BigDecimal getMinOrderAmount() {
        return minOrderAmount;
    }

    public void setMinOrderAmount(BigDecimal minOrderAmount) {
        this.minOrderAmount = minOrderAmount;
    }

    public BigDecimal getMaxAmountApplied() {
        return maxAmountApplied;
    }

    public void setMaxAmountApplied(BigDecimal maxAmountApplied) {
        this.maxAmountApplied = maxAmountApplied;
    }

    public boolean isActive() {
        return active;
    }

    public void setActive(boolean active) {
        this.active = active;
    }
    
    public Integer getTotalUsageCount() {
        return totalUsageCount;
    }

    public void setTotalUsageCount(Integer totalUsageCount) {
        this.totalUsageCount = totalUsageCount;
    }

    public BigDecimal getTotalSavingsAmount() {
        return totalSavingsAmount;
    }

    public void setTotalSavingsAmount(BigDecimal totalSavingsAmount) {
        this.totalSavingsAmount = totalSavingsAmount;
    }

    public Integer getUniqueUsersCount() {
        return uniqueUsersCount;
    }

    public void setUniqueUsersCount(Integer uniqueUsersCount) {
        this.uniqueUsersCount = uniqueUsersCount;
    }
}
