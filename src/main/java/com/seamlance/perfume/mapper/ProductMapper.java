package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.Product;
import com.seamlance.perfume.info.ProductInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = {ProductMapper.class})
public interface ProductMapper {

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "categoryInfo", source = "category")
    @Mapping(target = "discountInfo", source = "discount")
    ProductInfo toInfo(Product product);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "categoryInfo", target = "category")
    @Mapping(source = "discountInfo", target = "discount")
    Product toEntity(ProductInfo productInfo);

    List<ProductInfo> toConvertInfoList(List<Product> products);
}
