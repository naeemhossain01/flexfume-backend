package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.DiscountInfo;
import com.seamlance.perfume.mapper.DiscountMapper;
import com.seamlance.perfume.service.DiscountService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Map;


@Component
public class DiscountDelegateImpl implements DiscountDelegate {

    @Autowired
    private DiscountService discountService;

    @Autowired
    private DiscountMapper discountMapper;

    @Override
    public List<DiscountInfo> addDiscount(List<DiscountInfo> discountInfoList) {
        return discountMapper.toInfoList(discountService.addDiscount(discountMapper.toEntityList(discountInfoList)));
    }

    @Override
    public List<DiscountInfo> updateDiscount(Map<String, Integer> discountInfoList) {
        return discountMapper.toInfoList(discountService.updateDiscount(discountInfoList));
    }

    @Override
    public List<DiscountInfo> getAllProductDiscount() {
        return discountMapper.toInfoList(discountService.getAllProductDiscount());
    }

    @Override
    public void deleteDiscount(String id) {
        discountService.deleteDiscount(id);
    }
}
