package com.seamlance.perfume.mapper;


import com.seamlance.perfume.entity.Coupon;
import com.seamlance.perfume.info.CouponInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface CouponMapper {
    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    Coupon toEntity(CouponInfo couponInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    CouponInfo toInfo(Coupon coupon);

    List<CouponInfo> toInfoList(List<Coupon> couponList);
}
