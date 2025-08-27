package com.seamlance.perfume.service;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.DiscountDao;
import com.seamlance.perfume.entity.Discount;
import com.seamlance.perfume.exception.ResourceNotFoundException;

@Service
public class DiscountServiceImpl implements DiscountService {

    @Autowired
    private ProductService productService;

    @Autowired
    private DiscountDao discountDao;

    @Override
    public List<Discount> addDiscount(List<Discount> discounts) {
       try {
           discounts = discountDao.saveAllAndFlush(discounts);
       } catch (DataAccessException e) {
           e.printStackTrace();
       }

       return discounts;
    }

    @Override
    public List<Discount> updateDiscount(Map<String, Integer> discounts) {
        List<Discount> updatableDiscount = discountDao.findAll();

        if(updatableDiscount.isEmpty()) {
            throw new ResourceNotFoundException(ErrorConstant.DISCOUNT_NOT_FOUND);
        }

        for (Discount discount: updatableDiscount) {
            if(discounts.containsKey(discount.getId())) {
                discount.setPercentage(discounts.get(discount.getId()));
            }
        }

        try {
            updatableDiscount = discountDao.saveAll(updatableDiscount);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return updatableDiscount;
    }

    @Override
    public List<Discount> getAllProductDiscount() {
        List<Discount> discounts = new ArrayList<>();

        try {
            discounts = discountDao.findAll();
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return discounts;
    }

    @Override
    public void deleteDiscount(String id) {
        try {
            Discount discount = discountDao.findById(id)
                    .orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.DISCOUNT_NOT_FOUND));
            discountDao.delete(discount);
        } catch (DataAccessException e) {
            e.printStackTrace();
            throw new RuntimeException("Failed to delete discount", e);
        }
    }
}
