package com.seamlance.perfume.service;

import com.seamlance.perfume.entity.Category;

import java.util.List;

public interface CategoryService {
    Category createCategory(Category category);
    Category updateCategory(String id, Category updatedCategory);
    List<Category> getAllCategories();
    Category getCategoryById(String id);
    void deleteCategory(String id);
}
