SELECT
  t.*,
  dat.attributeID,
  dat.attributeName
FROM invTypes t
  INNER JOIN dgmTypeAttributes d ON d.typeID = t.typeID
  INNER JOIN dgmAttributeTypes a
    ON a.attributeID = d.attributeID AND a.attributeName IN ('primaryAttribute', 'secondaryAttribute')
  INNER JOIN dgmAttributeTypes dat ON (dat.attributeID = d.valueFloat OR dat.attributeID = d.valueInt)
WHERE t.groupID IN
      (
        SELECT g.groupID
        FROM invGroups g
        WHERE g.categoryID IN (
          SELECT categoryID
          FROM invCategories
          WHERE categoryName = 'Skill'
        )
      )
ORDER BY t.typeName, dat.attributeName;

--select distinct a.* from (
SELECT *
FROM invTypes t
  INNER JOIN dgmTypeAttributes d ON d.typeID = t.typeID
  INNER JOIN dgmAttributeTypes a
    ON a.attributeID = d.attributeID AND a.attributeName IN ('primaryAttribute', 'secondaryAttribute')
  INNER JOIN dgmAttributeTypes dat ON (dat.attributeID = d.valueFloat OR dat.attributeID = d.valueInt)
WHERE t.groupID IN
      (
        SELECT g.groupID
        FROM invGroups g
        WHERE g.categoryID IN (
          SELECT categoryID
          FROM invCategories
-- where categoryName = 'Skill'
        )
      )
      AND t.typeName IN ('Accounting')
ORDER BY t.typeName
--) as a order by a.typeName