package scalable.p2.thumbnail.repository;


import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

@Repository
public interface StatusRepository extends JpaRepository<Status, Integer> {

    @Query("select t from Status t where t.BucketID = :id and t.VideoName = :id2")
    Status getByBucketIDAndVideoName(@Param("id") Integer BucketID, @Param("id2") String VideoName);

    @Query("select t from Status t where t.BucketID = :id and t.Compiled = t.ExpectedFrames")
    List<Status> getCompletedGifs(@Param("id") Integer BucketID);
}
