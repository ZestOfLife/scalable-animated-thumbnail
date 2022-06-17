package scalable.p2.thumbnail.repository;

import lombok.*;
import org.hibernate.annotations.CreationTimestamp;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@Entity
@NoArgsConstructor
@AllArgsConstructor
@Builder
@Setter
@Getter
@IdClass(StatusID.class)
public class Status implements Serializable {
    @Id
    @Column(name = "bucket_id")
    private Integer BucketID;

    @Id
    @Column(name = "video_name")
    private String VideoName;

    @Column(name = "expected_frames")
    private Integer ExpectedFrames;

    @Column(name = "extracted")
    private Integer Extracted;

    @Column(name = "resized")
    private Integer Resized;

    @Column(name = "compiled")
    private Integer Compiled;
}
