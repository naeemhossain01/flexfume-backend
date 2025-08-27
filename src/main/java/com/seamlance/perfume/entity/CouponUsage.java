package com.seamlance.perfume.entity;

import com.seamlance.perfume.audit.AbstractAudit;
import jakarta.persistence.*;

import java.math.BigDecimal;

@Entity
@Table(name = "COUPON_USAGE")
@AttributeOverride(name = "id", column = @Column(name = "COUPON_USAGE_ID"))
public class CouponUsage extends AbstractAudit {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @ManyToOne
    @JoinColumn(name = "COUPON_ID")
    private Coupon coupon;

    @ManyToOne
    @JoinColumn(name = "USER_ID")
    private User user;

    @Column(name = "USAGE_COUNT")
    private int usageCount;

    @Column(name = "DISCOUNT_AMOUNT")
    private BigDecimal discountedAmount;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Coupon getCoupon() {
        return coupon;
    }

    public void setCoupon(Coupon coupon) {
        this.coupon = coupon;
    }

    public User getUser() {
        return user;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public int getUsageCount() {
        return usageCount;
    }

    public void setUsageCount(int usageCount) {
        this.usageCount = usageCount;
    }

    public BigDecimal getDiscountedAmount() {
        return discountedAmount;
    }

    public void setDiscountedAmount(BigDecimal discountedAmount) {
        this.discountedAmount = discountedAmount;
    }
}
