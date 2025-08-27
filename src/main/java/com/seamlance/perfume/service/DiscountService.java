package com.seamlance.perfume.service;

import java.util.List;
import java.util.Map;

import com.seamlance.perfume.entity.Discount;

public interface DiscountService {
    List<Discount> addDiscount(List<Discount> discounts);
    List<Discount> updateDiscount(Map<String, Integer> discounts);
    List<Discount> getAllProductDiscount();
    void deleteDiscount(String id);
}
