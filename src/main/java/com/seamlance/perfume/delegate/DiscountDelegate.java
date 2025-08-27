package com.seamlance.perfume.delegate;

import java.util.List;
import java.util.Map;

import com.seamlance.perfume.info.DiscountInfo;

public interface DiscountDelegate {
    List<DiscountInfo> addDiscount(List<DiscountInfo> discountInfoList);
    List<DiscountInfo> updateDiscount(Map<String, Integer> discountInfoList);
    List<DiscountInfo> getAllProductDiscount();
    void deleteDiscount(String id);
}
