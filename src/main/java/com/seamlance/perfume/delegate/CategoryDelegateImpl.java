package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.CategoryInfo;
import com.seamlance.perfume.mapper.CategoryMapper;
import com.seamlance.perfume.service.CategoryService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class CategoryDelegateImpl implements CategoryDelegate {

    @Autowired
    private CategoryService categoryService;

    @Autowired
    private CategoryMapper categoryMapper;

    @Override
    public CategoryInfo createCategory(CategoryInfo categoryInfo) {
        return categoryMapper.toInfo(categoryService.createCategory(categoryMapper.toEntity(categoryInfo)));
    }

    @Override
    public CategoryInfo updateCategory(String id, CategoryInfo categoryInfo) {
        return categoryMapper.toInfo(categoryService.updateCategory(id, categoryMapper.toEntity(categoryInfo)));
    }

    @Override
    public List<CategoryInfo> getAllCategories() {
        return categoryMapper.toInfoList(categoryService.getAllCategories());
    }

    @Override
    public CategoryInfo getCategoryById(String id) {
        return categoryMapper.toInfo(categoryService.getCategoryById(id));
    }

    @Override
    public void deleteCategory(String id) {
        categoryService.deleteCategory(id);
    }
}
