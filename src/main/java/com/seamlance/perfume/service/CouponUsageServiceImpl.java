package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.CouponUsageDao;
import com.seamlance.perfume.entity.*;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDateTime;
import java.util.List;

@Service
public class CouponUsageServiceImpl implements CouponUsageService {

    @Autowired
    private CouponService couponService;

    @Autowired
    private UserService userService;

    @Autowired
    private CouponUsageDao couponUsageDao;

    @Autowired
    private CartService cartService;

    private static final BigDecimal ONE_HUNDRED = new BigDecimal("100");

    @Override
    public BigDecimal applyCoupon(List<String> cartIdList, String couponCode) throws InvalidRequestsException {
        Coupon coupon = couponService.getCouponByCode(couponCode);
        User user = userService.getLoginUser();
        CouponUsage couponUsage = couponUsageDao.findByCouponCodeAndUserId(coupon.getCode(), user.getId());

        if(couponUsage != null && couponUsage.getUsageCount() >= coupon.getUsageLimit()) {
            throw new InvalidRequestsException("Coupon usage limit exceed");
        }

        if(!coupon.isActive() || coupon.getExpirationTime().isBefore(LocalDateTime.now())) {
            throw new InvalidRequestsException("Coupon is invalid");
        }

        List<Object> cartListObj = cartService.getCartAndProduct(cartIdList);

        BigDecimal productTotalAmount = calculateAmount(cartListObj);

        if(productTotalAmount.compareTo(coupon.getMinOrderAmount()) < 0) {
            throw new InvalidRequestsException("Please add more item");
        }

        BigDecimal couponAmount = calculateCouponAmount(coupon, productTotalAmount).min(coupon.getMaxAmountApplied());
        BigDecimal finalAmount = productTotalAmount.subtract(couponAmount);

        CouponUsage addCouponUsage = couponUsage == null ? addCouponUsage(coupon, user, couponAmount) : updateCouponUsage(couponUsage, couponAmount);

        try {
            addCouponUsage = couponUsageDao.save(addCouponUsage);
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return finalAmount;
    }

    @Override
    public CouponUsage getCouponUsage(String code, String userId) {
        CouponUsage couponUsage = couponUsageDao.findByCouponCodeAndUserId(code, userId);

        if(couponUsage == null) {
            throw new ResourceNotFoundException(ErrorConstant.NO_COUPON_USED);
        }

        return couponUsage;
    }

    @Override
    public void deleteCouponUsage(String couponCode) {
        User user = userService.getLoginUser();
        CouponUsage couponUsage = this.getCouponUsage(couponCode, user.getId());

        try {
            couponUsageDao.deleteById(couponUsage.getId());
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    public CouponUsage addCouponUsage(Coupon coupon, User user, BigDecimal couponAmount) {
        CouponUsage couponUsage = new CouponUsage();
        couponUsage.setCoupon(coupon);
        couponUsage.setUser(user);
        couponUsage.setUsageCount(1);
        couponUsage.setDiscountedAmount(couponAmount);

        return couponUsage;
    }

    public CouponUsage updateCouponUsage(CouponUsage couponUsage, BigDecimal couponAmount) {
        int usageLimit = couponUsage.getUsageCount();
        couponUsage.setUsageCount(usageLimit + 1);
        couponUsage.setDiscountedAmount(couponAmount);

        return couponUsage;
    }



    private BigDecimal calculateAmount(List<Object> cartListObj) {
        BigDecimal amount = BigDecimal.ZERO;

        for (Object obj : cartListObj) {
            Object[] data = (Object[]) obj;

            int quantity = data[1] != null ? (int) data[1] : 0;
            BigDecimal price = data[3] != null ? (BigDecimal) data[3] : null;
            int discount = data[4] != null ? (int) data[4] : 0;

            BigDecimal productPrice = calculateProductPrice(quantity, price, discount);
            amount = amount.add(productPrice);
        }

        return amount;
    }

    private BigDecimal calculateProductPrice(int quantity, BigDecimal price, int discount) {
        BigDecimal productIndividualPrice = discount > 0
                ? price.subtract(calculatePercentageAmount(price, discount))
                : price;

        return productIndividualPrice.multiply(BigDecimal.valueOf(quantity));
    }

    private BigDecimal calculatePercentageAmount(BigDecimal price, int percentage) {
        return price.multiply(BigDecimal.valueOf(percentage)).divide(ONE_HUNDRED, 2, RoundingMode.HALF_DOWN);
    }

    private BigDecimal calculateCouponAmount(Coupon coupon, BigDecimal amount) {
        BigDecimal calculatedAmount = new BigDecimal("0.0");

        if(coupon.getCouponType().equalsIgnoreCase("percentage")) {
            calculatedAmount = calculatePercentageAmount(amount, coupon.getAmount().intValue());
        }

        if(coupon.getCouponType().equalsIgnoreCase("fixed")) {
            calculatedAmount = coupon.getAmount();
        }

        return calculatedAmount;
    }
}
