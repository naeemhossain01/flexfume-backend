package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.ProductInfo;
import com.seamlance.perfume.mapper.ProductMapper;
import com.seamlance.perfume.service.ProductService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@Component
public class ProductDelegateImpl implements ProductDelegate{

    @Autowired
    private ProductMapper productMapper;

    @Autowired
    private ProductService productService;

    @Override
    public ProductInfo addProduct(ProductInfo productInfo) {
        return productMapper.toInfo(
                productService.addProduct(productMapper.toEntity(productInfo)));
    }

    @Override
    public ProductInfo getProductById(String id) {
        return productMapper.toInfo(productService.getProductById(id));
    }

    @Override
    public ProductInfo updateProduct(String id, ProductInfo productInfo) {
        return productMapper.toInfo(
                productService.updateProduct(id, productMapper.toEntity(productInfo)));
    }

    @Override
    public List<ProductInfo> getAllProduct() {
        return productMapper.toConvertInfoList(productService.getAllProduct());
    }

    @Override
    public void deleteProduct(String id) {
        productService.deleteProduct(id);
    }

    @Override
    public List<ProductInfo> getProductByCategory(String categoryId) {
        return productMapper.toConvertInfoList(productService.getProductByCategory(categoryId));
    }

    @Override
    public List<ProductInfo> searchProduct(String value) {
        return productMapper.toConvertInfoList(productService.searchProduct(value));
    }

    @Override
    public ProductInfo uploadImage(String productId, MultipartFile file) {
        return productMapper.toInfo(productService.uploadImage(productId, file));
    }
}
