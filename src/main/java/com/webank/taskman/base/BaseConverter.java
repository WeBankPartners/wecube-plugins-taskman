package com.webank.taskman.base;

import java.util.List;

public interface BaseConverter<D, E> {

    /**
     * DTO to Entity
     * @param dto
     * @return
     */
    E toEntity(D dto);

    /**
     * Entity to DTO
     * @param entity
     * @return
     */
    /*@Mappings({
            @Mapping(target ="createBy",ignore = true),
            @Mapping(target ="createTime",ignore = true),
            @Mapping(target ="updateBy",ignore = true),
            @Mapping(target ="updateTime",ignore = true),
            @Mapping(target ="orderBy",ignore = true),
            @Mapping(target ="optimisticLock",ignore = true),
            @Mapping(target ="params",ignore = true),
            @Mapping(target ="ids",ignore = true),
            @Mapping(target ="search",ignore = true),
            @Mapping(target ="searchSql",ignore = true),
    })*/
    D toDto(E entity);

    /**
     * DTOList  to EntityList
     * @param dtoList
     * @return
     */
    List<E> toEntity(List<D> dtoList);

    /**
     * EntityList to  DTOList
     * @param entityList
     * @return
     */
    List <D> toDto(List<E> entityList);
}
