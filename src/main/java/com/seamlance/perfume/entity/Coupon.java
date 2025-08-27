package com.seamlance.perfume.entity;

import com.fasterxml.jackson.annotation.JsonManagedReference;
import com.seamlance.perfume.audit.AbstractAudit;
import jakarta.persistence.*;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;

@Entity
@Table(name = "COUPON")
@AttributeOverride(name = "id", column = @Column(name = "COUPON_ID"))
public class Coupon extends AbstractAudit {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "COUPON_TYPE")
    private String couponType;

    @Column(name = "CODE")
    private String code;

    @Column(name = "EXPIRATION_TIME")
    private LocalDateTime expirationTime;

    @Column(name = "USAGE_LIMIT")
    private int usageLimit;

    @Column(name = "AMOUNT")
    private BigDecimal amount;

    @Column(name = "MIN_ORDER_AMOUNT")
    private BigDecimal minOrderAmount;

    @Column(name = "MAX_AMOUNT_APPLIED")
    private BigDecimal maxAmountApplied;

    @Column(name = "IS_ACTIVE")
    private boolean active;

    @OneToMany(mappedBy = "coupon", fetch = FetchType.LAZY, cascade = CascadeType.ALL, orphanRemoval = true)
    @JsonManagedReference
    private List<CouponUsage> couponUsage;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

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

    public LocalDateTime getExpirationTime() {
        return expirationTime;
    }

    public void setExpirationTime(LocalDateTime expirationTime) {
        this.expirationTime = expirationTime;
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

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public int getUsageLimit() {
        return usageLimit;
    }

    public void setUsageLimit(int usageLimit) {
        this.usageLimit = usageLimit;
    }

    public boolean isActive() {
        return active;
    }

    public void setActive(boolean active) {
        this.active = active;
    }
}
