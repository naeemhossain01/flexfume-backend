package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.CouponUsage;
import com.seamlance.perfume.info.CouponUsageInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring", uses = {CouponMapper.class, UserMapper.class})
public interface CouponUsageMapper {
    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    CouponUsage toEntity(CouponUsageInfo couponUsageInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    CouponUsageInfo toInfo(CouponUsage couponUsage);
}
