package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.Discount;
import com.seamlance.perfume.info.DiscountInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = {ProductMapper.class})
public interface DiscountMapper {
    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "productInfo", target = "product")
    Discount toEntity(DiscountInfo discountInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "productInfo", source = "product")
    DiscountInfo toInfo(Discount discount);

    List<DiscountInfo> toInfoList(List<Discount> discounts);

    List<Discount> toEntityList(List<DiscountInfo> discountInfos);
}
