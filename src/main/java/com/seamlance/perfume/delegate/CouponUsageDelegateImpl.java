package com.seamlance.perfume.delegate;

import com.seamlance.perfume.entity.Cart;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.info.CartInfo;
import com.seamlance.perfume.mapper.CartMapper;
import com.seamlance.perfume.service.CouponUsageService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.util.List;

@Component
public class CouponUsageDelegateImpl implements CouponUsageDelegate {

    @Autowired
    private CartMapper cartMapper;

    @Autowired
    private CouponUsageService couponUsageService;

    @Override
    public BigDecimal applyCoupon(List<String> cartInfoList, String couponCode) throws InvalidRequestsException {

        return couponUsageService.applyCoupon(cartInfoList, couponCode);
    }

    @Override
    public void deleteCouponUsage(String code) {
        couponUsageService.deleteCouponUsage(code);
    }
}
