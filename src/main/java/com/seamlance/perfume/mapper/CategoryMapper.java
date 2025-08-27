package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.Category;
import com.seamlance.perfume.info.CategoryInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = {ProductMapper.class})
public interface CategoryMapper {
    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "productInfos", target = "productList")
    Category toEntity(CategoryInfo categoryInfo);


    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "productInfos", source = "productList")
    CategoryInfo toInfo(Category category);

    List<CategoryInfo> toInfoList(List<Category> categoryList);
}
