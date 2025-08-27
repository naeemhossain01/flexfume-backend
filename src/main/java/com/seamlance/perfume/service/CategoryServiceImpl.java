package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.CategoryDao;
import com.seamlance.perfume.entity.Category;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

@Service
public class CategoryServiceImpl implements CategoryService {

    @Autowired
    private CategoryDao categoryDao;

    @Override
    public Category createCategory(Category category) {
        //TODO: validation need

        try {
            category = categoryDao.saveAndFlush(category);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return category;
    }

    @Override
    public Category updateCategory(String id, Category updatedCategory) {
        Category category = categoryDao.findById(id).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.CATEGORY_NOT_FOUND));

        try {
            category.setName(updatedCategory.getName());
            category = categoryDao.save(category);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return category;
    }

    @Override
    public List<Category> getAllCategories() {
        List<Category> categories = new ArrayList<>();

        try {
            categories = categoryDao.findAll();
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return categories;
    }

    @Override
    public Category getCategoryById(String id) {
        return categoryDao.findById(id).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.CATEGORY_NOT_FOUND));
    }

    @Override
    public void deleteCategory(String id) {
        categoryDao.findById(id).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.CATEGORY_NOT_FOUND));

        try {
            categoryDao.deleteById(id);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }
    }
}
