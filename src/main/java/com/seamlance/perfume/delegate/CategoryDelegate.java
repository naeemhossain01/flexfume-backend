package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.CategoryInfo;

import java.util.List;

public interface CategoryDelegate {
    CategoryInfo createCategory(CategoryInfo categoryInfo);
    CategoryInfo updateCategory(String id, CategoryInfo categoryInfo);
    List<CategoryInfo> getAllCategories();
    CategoryInfo getCategoryById(String id);
    void deleteCategory(String id);
}
